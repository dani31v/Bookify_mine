package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"bookify/pkg/health"
	"bookify/pkg/httpserver"

	"bookify/services/shelves/internal/controller"
	httpH "bookify/services/shelves/internal/handler/http"
	"bookify/services/shelves/internal/repository/memory"
	pgRepo "bookify/services/shelves/internal/repository/postgres"

	shelvesgrpc "bookify/proto/shelves_proto/transport/grpc"
	grpcTransport "bookify/services/shelves/internal/transport/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	serviceName = "shelves"
	httpPort    = 8084
	grpcPort    = 50055
)

func main() {

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, pass, dbname, port,
	)

	var repo controller.Repository
	pg, err := pgRepo.New(dsn)
	if err != nil {
		log.Printf("ERROR connecting to Postgres: %v", err)
		log.Println("Falling back to memory repository")

		repo = memory.New()
	} else {
		log.Println("Using Postgres repository for shelves")
		repo = pg
	}

	ctrl := controller.NewShelves(repo)

	go func() {
		lis, err := net.Listen("tcp", ":50055")
		if err != nil {
			log.Fatalf("Failed to listen on gRPC port: %v", err)
		}

		grpcServer := grpc.NewServer()
		reflection.Register(grpcServer)

		shelvesgrpc.RegisterShelvesServiceServer(
			grpcServer,
			grpcTransport.NewShelvesServer(ctrl),
		)

		log.Printf("Shelves gRPC server running on port %d", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	httpserver.Start(httpPort, func(mux *http.ServeMux) {
		mux.HandleFunc("/healthz", health.Handler)
		mux.HandleFunc("/shelf", httpH.ShelfHandler(ctrl))
	})
}
