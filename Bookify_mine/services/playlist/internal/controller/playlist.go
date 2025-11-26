package controller

import (
	"errors"

	"bookify/services/playlist/internal/model"
)

// Interface común para evitar depender de memory/postgres directamente.
type Repository interface {
	ForBook(bookID string) model.Playlist
	Create(playlist model.Playlist) error
}

var (
	ErrAlreadyExists  = errors.New("playlist already exists")
	ErrBookIDRequired = errors.New("bookId required")
)

type Playlists struct {
	repo Repository
}

func NewPlaylists(repo Repository) *Playlists {
	return &Playlists{repo: repo}
}

// Devuelve SIEMPRE un model.Playlist (incluso si está vacío)
func (c *Playlists) ForBook(id string) model.Playlist {
	return c.repo.ForBook(id)
}

func (c *Playlists) Create(playlist model.Playlist) (model.Playlist, error) {
	if playlist.BookID == "" {
		return model.Playlist{}, ErrBookIDRequired
	}

	// Llama al repo correcto (Postgres o Memory)
	if err := c.repo.Create(playlist); err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			return model.Playlist{}, ErrAlreadyExists
		}
		return model.Playlist{}, err
	}

	return playlist, nil
}
