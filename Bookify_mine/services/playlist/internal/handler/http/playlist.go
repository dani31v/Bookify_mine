package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"bookify/services/playlist/internal/controller"
	"bookify/services/playlist/internal/model"
)

func PlaylistHandler(c *controller.Playlists) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id := r.URL.Query().Get("bookId")
			if id == "" {
				http.Error(w, "bookId required", http.StatusBadRequest)
				return
			}

			_ = json.NewEncoder(w).Encode(c.ForBook(id))

		case http.MethodPost:
			var playlist model.Playlist
			if err := json.NewDecoder(r.Body).Decode(&playlist); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			created, err := c.Create(playlist)
			if err != nil {
				switch {
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
			_ = json.NewEncoder(w).Encode(created)

		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
