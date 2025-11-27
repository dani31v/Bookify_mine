package rest

import (
	"encoding/json"
	"net/http"

	c "bookify/services/gateway/internal/clients"

	librarygrpc "bookify/proto/library_proto/transport/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LibraryHandler struct {
	Library *c.LibraryClient
}

func NewLibraryHandler(lib *c.LibraryClient) *LibraryHandler {
	return &LibraryHandler{Library: lib}
}

type bookPayload struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Pages   int32  `json:"pages"`
	Edition string `json:"edition"`
}

func (h *LibraryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleCreateBook(w, r)
	default:
		w.Header().Set("Allow", "POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *LibraryHandler) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var payload bookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid json payload", http.StatusBadRequest)
		return
	}

	book := &librarygrpc.Book{
		Id:      payload.ID,
		Title:   payload.Title,
		Author:  payload.Author,
		Pages:   payload.Pages,
		Edition: payload.Edition,
	}

	created, err := h.Library.CreateBook(r.Context(), book)
	if err != nil {
		handleLibraryError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(created)
}

func handleLibraryError(w http.ResponseWriter, err error) {
	if st, ok := status.FromError(err); ok {
		code := http.StatusInternalServerError
		switch st.Code() {
		case codes.InvalidArgument:
			code = http.StatusBadRequest
		case codes.AlreadyExists:
			code = http.StatusConflict
		}
		http.Error(w, st.Message(), code)
		return
	}

	http.Error(w, "library service error: "+err.Error(), http.StatusInternalServerError)
}
