package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/caarlos0/env/v10"
	_ "github.com/glebarez/go-sqlite"
	"github.com/gorilla/mux"
	"github.com/thoom/mimic/acme"
	"github.com/thoom/mimic/mimic"
)

func main() {
	// Parse environment variables
	cfg := mimic.Config{}
	options := env.Options{
		OnSet: func(tag string, value interface{}, isDefault bool) {
			if tag == "HOST_URL" {
				cfg.HostURL = value.(string)
				url, _ := url.Parse(cfg.HostURL)
				cfg.ParsedHostName.Protocol = url.Scheme
				cfg.ParsedHostName.Host = url.Hostname()
				cfg.ParsedHostName.Port = url.Port()
			}
		},
	}

	if err := env.ParseWithOptions(&cfg, options); err != nil {
		fmt.Printf("%+v\n", err)
	}

	// Connect to DB
	db, err := sql.Open("sqlite", "mimic.db?_journal_mode=WAL")
	if err != nil {
		log.Fatal(err)
	}

	// Create DB schema if it doesn't exist
	mdb := mimic.MimicDB{DB: db}
	mdb.CreateSchema()

	// Set up routes
	r := mux.NewRouter()
	r.Use(mimic.LoggingMiddleware)
	r.Use(func(next http.Handler) http.Handler {
		return mimic.JoseMiddleware(next, &mdb)
	})

	r.HandleFunc("/directory", func(w http.ResponseWriter, r *http.Request) {
		acme.DirectoryHandler(w, r, &cfg)
	})

	r.HandleFunc("/terms/v1", mimic.TermsHandler).Methods("GET")

	r.PathPrefix("/acme/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// no op - will be handled by runsAfter function
	}).Methods("HEAD")

	r.PathPrefix("/acme/").HandlerFunc(acme.GenericHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.ParsedHostName.Port), acme.AddNonceToResponse(r, &mdb)))
}
