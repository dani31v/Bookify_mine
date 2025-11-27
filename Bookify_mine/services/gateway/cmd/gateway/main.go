package main

import (
	"log"
	"net/http"

	"bookify/pkg/httpserver"

	c "bookify/services/gateway/internal/clients"
	gatewaygrpc "bookify/services/gateway/internal/handler/grpc"
	resthandler "bookify/services/gateway/internal/handler/rest"
)

const serviceName = "gateway"

func main() {
	port := 8090

	lib, err := c.NewLibraryClient("library:50051")
	if err != nil {
		log.Fatal("cannot create library client:", err)
	}

	rev, err := c.NewReviewsClient("reviews:50052")
	if err != nil {
		log.Fatal("cannot create reviews client:", err)
	}

	pl, err := c.NewPlaylistClient("playlist:50053")
	if err != nil {
		log.Fatal("cannot create playlist client:", err)
	}

	sh, err := c.NewShelvesClient("shelves:50055")
	if err != nil {
		log.Fatal("cannot create shelves client:", err)
	}

	h := gatewaygrpc.NewOverviewHandler(lib, rev, pl, sh)
	libraryHandler := resthandler.NewLibraryHandler(lib)

	httpserver.Start(port, func(mux *http.ServeMux) {
		*mux = *gatewaygrpc.NewMux(h, libraryHandler)
	})
}
