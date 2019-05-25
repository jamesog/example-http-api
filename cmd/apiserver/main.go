package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jamesog/example-http-api/api"
	"github.com/jamesog/example-http-api/database/postgres"
)

func main() {
	listen := os.Getenv("LISTEN_ADDR")
	if listen == "" {
		listen = ":8000"
	}

	db, err := postgres.NewDB("sslmode=disable user=postgres dbname=example")
	if err != nil {
		log.Fatal(err)
	}
	apisvc := api.NewService(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/", apisvc.Routes())

	log.Fatal(http.ListenAndServe(listen, r))
}
