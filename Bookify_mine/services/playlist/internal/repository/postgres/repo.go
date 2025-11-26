package postgres

import (
	"errors"
	"fmt"

	"bookify/services/playlist/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrPlaylistExists = errors.New("playlist already exists")
	ErrBookIDRequired = errors.New("bookId required")
)

type Repo struct {
	db *gorm.DB
}

func New(dsn string) (*Repo, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgres connect error: %w", err)
	}

	// Create table if not exists
	if err := db.AutoMigrate(&model.Playlist{}); err != nil {
		return nil, fmt.Errorf("migrate error: %w", err)
	}

	return &Repo{db: db}, nil
}

// ------------------------------------------
// Implements Repository interface
// ------------------------------------------

// Get playlist for a book
func (r *Repo) ForBook(bookID string) model.Playlist {
	var playlist model.Playlist
	err := r.db.First(&playlist, "book_id = ?", bookID).Error

	// If not found, return empty playlist with only BookID
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Playlist{BookID: bookID}
	}

	return playlist
}

// Create playlist
func (r *Repo) Create(p model.Playlist) error {
	if p.BookID == "" {
		return ErrBookIDRequired
	}

	// ID is required as primary key, so we generate it
	p.ID = "playlist:" + p.BookID

	// Prevent duplicates
	var existing model.Playlist
	err := r.db.First(&existing, "id = ?", p.ID).Error
	if err == nil {
		return ErrPlaylistExists
	}

	return r.db.Create(&p).Error
}
