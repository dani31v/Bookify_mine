package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"bookify/services/reviews/internal/controller"
	"bookify/services/reviews/internal/model"
)

func ReviewsHandler(c *controller.Reviews) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id := r.URL.Query().Get("bookId")
			if id == "" {
				http.Error(w, "bookId required", http.StatusBadRequest)
				return
			}

			_ = json.NewEncoder(w).Encode(map[string]any{"data": c.ByBook(id)})

		case http.MethodPost:
			var review model.Review
			if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			created, err := c.Create(review)
			if err != nil {
				switch {
				case errors.Is(err, controller.ErrIDRequired):
					http.Error(w, err.Error(), http.StatusBadRequest)
				case errors.Is(err, controller.ErrBookIDRequired):
					http.Error(w, err.Error(), http.StatusBadRequest)
				case errors.Is(err, controller.ErrAlreadyExists):
					http.Error(w, err.Error(), http.StatusConflict)
				default:
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
				return
			}

			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{"data": created})

		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
