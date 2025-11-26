package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"bookify/pkg/health"
	"bookify/pkg/httpserver"

	"bookify/services/reviews/internal/controller"
	httpH "bookify/services/reviews/internal/handler/http"
	memRepo "bookify/services/reviews/internal/repository/memory"
	pgRepo "bookify/services/reviews/internal/repository/postgres"

	reviewsgrpc "bookify/proto/reviews_proto/transport/grpc"
	grpcTransport "bookify/services/reviews/internal/transport/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	serviceName = "reviews"
	httpPort    = 8082
	grpcPort    = 50052
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
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, dbname, port,
	)

	var repo controller.Repository

	pg, err := pgRepo.New(dsn)
	if err != nil {
		log.Printf("ERROR using Postgres repo: %v — falling back to memory repo", err)
		repo = memRepo.New()
	} else {
		log.Println("Using Postgres repository for reviews")
		repo = pg
	}

	ctrl := controller.NewReviews(repo)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
		if err != nil {
			log.Fatalf("Failed to listen on gRPC port: %v", err)
		}

		grpcServer := grpc.NewServer()
		reflection.Register(grpcServer)

		reviewsgrpc.RegisterReviewServiceServer(
			grpcServer,
			grpcTransport.NewReviewsServer(ctrl),
		)

		log.Printf("Reviews gRPC server running on port %d", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	httpserver.Start(httpPort, func(mux *http.ServeMux) {
		mux.HandleFunc("/healthz", health.Handler)
		mux.HandleFunc("/reviews", httpH.ReviewsHandler(ctrl))
	})
}
