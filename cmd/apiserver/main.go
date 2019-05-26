package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jamesog/example-http-api/api"
	"github.com/jamesog/example-http-api/database"
	"github.com/jamesog/example-http-api/database/postgres"
	"google.golang.org/grpc"
)

func httpServer(db database.Storage) {
	listen := os.Getenv("HTTP_ADDR")
	if listen == "" {
		listen = ":8000"
	}

	apisvc := api.NewService(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/", apisvc.Routes())

	log.Fatal(http.ListenAndServe(listen, r))
}

func grpcServer(db database.Storage) {
	listen := os.Getenv("GRPC_ADDR")
	if listen == "" {
		listen = ":8001"
	}

	lis, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatal(err)
	}

	apisvc := api.NewService(db)

	grpcServer := grpc.NewServer()
	api.RegisterExampleServiceServer(grpcServer, apisvc)
	log.Fatal(grpcServer.Serve(lis))
}

func main() {
	db, err := postgres.NewDB("sslmode=disable user=postgres dbname=example")
	if err != nil {
		log.Fatal(err)
	}

	go httpServer(db)
	grpcServer(db)
}
