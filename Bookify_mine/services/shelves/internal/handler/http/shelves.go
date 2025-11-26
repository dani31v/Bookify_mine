package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"bookify/services/shelves/internal/controller"
	"bookify/services/shelves/internal/model"
)

func ShelfHandler(c *controller.Shelves) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			user := r.URL.Query().Get("userId")
			book := r.URL.Query().Get("bookId")
			if user == "" || book == "" {
				http.Error(w, "bookId required", http.StatusBadRequest)
				return
			}

			item, ok := c.Get(user, book)
			if !ok {
				http.NotFound(w, r)
				return
			}

			_ = json.NewEncoder(w).Encode(item)

		case http.MethodPost:
			var item model.ShelfItem
			if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			created, err := c.Create(item)
			if err != nil {
				switch {
				case errors.Is(err, controller.ErrUserIDRequired), errors.Is(err, controller.ErrBookIDRequired):
					http.Error(w, err.Error(), http.StatusBadRequest)
				case errors.Is(err, controller.ErrAlreadyExists):
					http.Error(w, err.Error(), http.StatusConflict)
				default:
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
				return
			}

			_ = json.NewEncoder(w).Encode(created)

		}
	}
}
