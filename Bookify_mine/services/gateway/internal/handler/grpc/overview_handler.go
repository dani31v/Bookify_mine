package grpc

import (
	"context"
	"encoding/json"
	"net/http"

	c "bookify/services/gateway/internal/clients"
)

type OverviewHandler struct {
	Library  *c.LibraryClient
	Reviews  *c.ReviewsClient
	Playlist *c.PlaylistClient
	Shelves  *c.ShelvesClient
}

func NewOverviewHandler(
	lib *c.LibraryClient,
	rev *c.ReviewsClient,
	pl *c.PlaylistClient,
	sh *c.ShelvesClient,
) *OverviewHandler {
	return &OverviewHandler{
		Library:  lib,
		Reviews:  rev,
		Playlist: pl,
		Shelves:  sh,
	}
}

func (h *OverviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	bookID := r.URL.Query().Get("bookId")
	userID := r.URL.Query().Get("userId")

	bookResp, err := h.Library.GetBook(ctx, bookID)
	if err != nil {
		http.Error(w, "library error: "+err.Error(), 500)
		return
	}

	reviewsResp, err := h.Reviews.GetReviews(ctx, bookID)
	if err != nil {
		http.Error(w, "reviews error: "+err.Error(), 500)
		return
	}

	playlistResp, err := h.Playlist.GetPlaylistForBook(ctx, bookID)
	if err != nil {
		http.Error(w, "playlist error: "+err.Error(), 500)
		return
	}

	var shelfResp any
	if userID != "" {
		s, err := h.Shelves.GetShelfItem(ctx, userID, bookID)
		if err == nil {
			shelfResp = s
		}
	}

	out := map[string]any{
		"book":     bookResp,
		"reviews":  reviewsResp,
		"playlist": playlistResp,
		"shelf":    shelfResp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}
