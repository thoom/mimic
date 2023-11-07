package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/caarlos0/env/v10"
	_ "github.com/glebarez/go-sqlite"
	"github.com/gorilla/mux"
	"github.com/thoom/mimic/acme"
	"github.com/thoom/mimic/mimic"
)

type MimicContext string

const (
	MimicJose MimicContext = "mimicJose"
)

func randString(length int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]rune, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func main() {
	cfg := mimic.Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", cfg)

	// Connect to DB
	db, err := sql.Open("sqlite", "mimic.db?_journal_mode=WAL")
	if err != nil {
		log.Fatal(err)
	}

	// Create DB schema if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS nonce (
		nonce TEXT PRIMARY KEY, 
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP, 
		used_at DATETIME DEFAULT NULL
		) WITHOUT ROWID;
`)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Use(joseMiddleware)

	r.HandleFunc("/directory", func(w http.ResponseWriter, r *http.Request) {
		acme.DirectoryHandler(w, r, &cfg)
	})

	r.PathPrefix("/acme/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonce := randString(32)
		db.Exec("INSERT INTO nonce (nonce) VALUES (?)", nonce)

		w.Header().Set("Replay-Nonce", nonce)
		w.Header().Set("Cache-Control", "no-store")
	}).Methods("HEAD")

	r.PathPrefix("/acme/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// no op for now

		joseJson := r.Context().Value(MimicJose).(mimic.JoseJson)
		if joseJson.Payload != "" {
			log.Printf("%+v\n", joseJson)
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add some route logging
		log.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Header.Get("Content-Type"))
		next.ServeHTTP(w, r)
	})
}

func joseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)

		var joseJson mimic.JoseJson
		if r.Header.Get("Content-Type") == "application/jose+json" {
			json.Unmarshal([]byte(buf.Bytes()), &joseJson)
			joseJson.DecodeProtected()
			joseJson.DecodePayload()

			//TODO: validate the JWT
		}

		ctx := context.WithValue(r.Context(), MimicJose, joseJson)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
