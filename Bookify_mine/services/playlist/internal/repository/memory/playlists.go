package memory

import (
	"errors"
	"sync"

	"bookify/services/playlist/internal/model"
)

var (
	ErrPlaylistExists = errors.New("playlist already exists")
	ErrPlaylistBookID = errors.New("bookId required")
)

type Repo struct {
	mu   sync.RWMutex
	data map[string]model.Playlist
}

func New() *Repo {
	return &Repo{data: map[string]model.Playlist{
		"book-1": {
			BookID: "book-1",
			Tracks: []model.Song{
				{Title: "Midnight Rain", Artist: "Taylor Swift"},
				{Title: "Perfect", Artist: "Ed Sheeran"},
			},
		},
	}}
}

func (r *Repo) ForBook(id string) model.Playlist {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if playlist, ok := r.data[id]; ok {
		return playlist
	}
	return model.Playlist{BookID: id}
}

func (r *Repo) Create(playlist model.Playlist) error {
	if playlist.BookID == "" {
		return ErrPlaylistBookID
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[playlist.BookID]; exists {
		return ErrPlaylistExists
	}

	r.data[playlist.BookID] = playlist
	return nil
}
