package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"bookify/services/library/internal/controller"
	"bookify/services/library/internal/model"
)

func BooksHandler(c *controller.Books) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id == "" {
				http.Error(w, "id required", http.StatusBadRequest)
				return
			}

			book, ok := c.Get(id)
			if !ok {
				http.NotFound(w, r)
				return
			}

			_ = json.NewEncoder(w).Encode(book)

		case http.MethodPost:
			var book model.Book
			if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			created, err := c.Create(book)
			if err != nil {
				switch {
				case errors.Is(err, controller.ErrIDRequired):
					http.Error(w, err.Error(), http.StatusBadRequest)
				case errors.Is(err, controller.ErrAlreadyExists):
					http.Error(w, err.Error(), http.StatusConflict)
				default:
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
				return
			}

			w.Header().Set("-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(created)

		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
