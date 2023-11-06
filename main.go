package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env/v10"
	_ "github.com/glebarez/go-sqlite"
	"github.com/gorilla/mux"
	"github.com/thoom/mimic/acme"
	"github.com/thoom/mimic/mimic"
)

func main() {
	cfg := mimic.Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", cfg)

	// connect
	db, err := sql.Open("sqlite", "test.db?_pragma=journal_mode(WAL)")
	if err != nil {
		log.Fatal(err)
	}

	// get SQLite version
	var v string
	db.QueryRow("select sqlite_version()").Scan(&v)
	fmt.Printf("v: %v\n", v)

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/discovery", func(w http.ResponseWriter, r *http.Request) {
		acme.DiscoveryHandler(w, r, &cfg)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
