package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jamesog/example-http-api/api"
	"github.com/jamesog/example-http-api/database"
	"github.com/jamesog/example-http-api/database/postgres"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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

type rpcapi struct {
	*api.API
	log zerolog.Logger
	mw  []grpc.UnaryServerInterceptor
}

func (api *rpcapi) use(mw func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)) {
	api.mw = append(api.mw, mw)
}

func (api *rpcapi) rpclog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	p, _ := peer.FromContext(ctx)
	remote, _, _ := net.SplitHostPort(p.Addr.String())

	start := time.Now()
	resp, err := handler(ctx, req)
	dur := time.Since(start)
	s, _ := status.FromError(err)
	api.log.Info().
		Str("rpc", info.FullMethod).
		Str("remote", remote).
		Str("code", s.Code().String()).
		Str("duration", dur.String()).
		Msg(s.Message())
	return resp, err
}

func (api *rpcapi) rpcmetric(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Pretend this records metrics
	start := time.Now()
	resp, err := handler(ctx, req)
	dur := time.Since(start)
	fmt.Printf("Metric: %s\n", dur)

	return resp, err
}

func (api *rpcapi) rpcMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if len(api.mw) == 0 {
		return handler(ctx, req)
	}

	resp, err := api.mw[len(api.mw)-1](ctx, req, info, handler)
	for i := len(api.mw) - 2; i >= 0; i-- {
		resp, err = api.mw[i](ctx, req, info, handler)
	}
	return resp, err
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

	l := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	apisvc := rpcapi{API: api.NewService(db), log: l}

	apisvc.use(apisvc.rpclog)
	apisvc.use(apisvc.rpcmetric)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(apisvc.mw...)))
	api.RegisterExampleServiceServer(grpcServer, apisvc)
	reflection.Register(grpcServer)
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
