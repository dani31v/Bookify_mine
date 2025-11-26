package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"bookify/pkg/health"
	"bookify/pkg/httpserver"

	"bookify/services/playlist/internal/controller"
	httpH "bookify/services/playlist/internal/handler/http"
	memoryRepo "bookify/services/playlist/internal/repository/memory"
	postgresRepo "bookify/services/playlist/internal/repository/postgres"
	grpcTransport "bookify/services/playlist/internal/transport/grpc"

	playlistgrpc "bookify/proto/playlist_proto/transport/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	serviceName = "playlist"
	httpPort    = 8083
	grpcPort    = 50053
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

	// Define repository interface (postgres or memory)
	var repo controller.Repository

	// Try Postgres first
	if host != "" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			host, user, pass, dbname, port,
		)

		pgRepo, err := postgresRepo.New(dsn)
		if err != nil {
			log.Printf("ERROR connecting to Postgres playlist repo: %v — using memory instead", err)
			repo = memoryRepo.New()
		} else {
			log.Println("Using Postgres repository for playlist")
			repo = pgRepo
		}
	} else {
		log.Println("DB_HOST not set — using memory playlist repo")
		repo = memoryRepo.New()
	}

	// ======== CONTROLLER ========
	ctrl := controller.NewPlaylists(repo)

	// ======== START GPRC ========
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
		if err != nil {
			log.Fatalf("Failed to listen on gRPC port: %v", err)
		}

		grpcServer := grpc.NewServer()
		reflection.Register(grpcServer)

		playlistgrpc.RegisterPlaylistServiceServer(grpcServer, grpcTransport.NewPlaylistServer(ctrl))

		log.Printf("Playlist gRPC server running on port %d", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// ======== START HTTP ========
	httpserver.Start(httpPort, func(mux *http.ServeMux) {
		mux.HandleFunc("/healthz", health.Handler)
		mux.HandleFunc("/playlist", httpH.PlaylistHandler(ctrl))
	})
}
