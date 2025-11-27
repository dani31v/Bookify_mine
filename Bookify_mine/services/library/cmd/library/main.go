package main

import (
	"log"
	"net"
	"net/http"

	"bookify/pkg/health"
	"bookify/pkg/httpserver"

	"bookify/services/library/internal/controller"
	httpH "bookify/services/library/internal/handler/http"
	memoryRepo "bookify/services/library/internal/repository/memory"
	postgresRepo "bookify/services/library/internal/repository/postgres"

	librarygrpc "bookify/proto/library_proto/transport/grpc"
	grpcTransport "bookify/services/library/internal/transport/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	serviceName = "library"
	httpPort    = 8081
	grpcPort    = 50051
)

func main() {

	dsn := "host=postgres user=userbookify password=user_bookify dbname=bookify port=5432 sslmode=disable"

	postgresRepoInstance, err := postgresRepo.New(dsn)
	if err != nil {
		log.Printf("WARN: postgres failed, falling back to memory: %v", err)
		postgresRepoInstance = nil
	}

	var repo controller.Repository

	if postgresRepoInstance != nil {
		repo = postgresRepoInstance
		log.Println("Using Postgres repository for library")
	} else {

		repo = memoryRepo.New()
		log.Println("Using in-memory repository for library")
	}

	ctrl := controller.NewBooks(repo)

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen on gRPC port: %v", err)
		}

		grpcServer := grpc.NewServer()
		reflection.Register(grpcServer)

		librarygrpc.RegisterLibraryServiceServer(grpcServer, grpcTransport.NewLibraryServer(ctrl))

		log.Printf("gRPC server running on port %d", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()
	httpserver.Start(httpPort, func(mux *http.ServeMux) {
		mux.HandleFunc("/healthz", health.Handler)
		mux.HandleFunc("/books", httpH.BooksHandler(ctrl))
	})
}
