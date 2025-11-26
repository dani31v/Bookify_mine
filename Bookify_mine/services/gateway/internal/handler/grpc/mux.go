package grpc

import (
	"net/http"
)

func NewMux(h http.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/overview/book", h)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	return mux
}
