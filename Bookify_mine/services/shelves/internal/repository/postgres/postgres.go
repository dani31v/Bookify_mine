package postgres

import (
	"errors"
	"fmt"

	"bookify/services/shelves/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrShelfExists    = errors.New("shelf item already exists")
	ErrUserIDRequired = errors.New("userId required")
	ErrBookIDRequired = errors.New("bookId required")
)

type Repo struct {
	db *gorm.DB
}

func New(dsn string) (*Repo, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Create table if not exists
	if err := db.AutoMigrate(&model.ShelfItem{}); err != nil {
		return nil, fmt.Errorf("failed to migrate shelves table: %w", err)
	}

	return &Repo{db: db}, nil
}

// Get a shelf item
func (r *Repo) Get(userID, bookID string) (model.ShelfItem, bool) {
	id := userID + ":" + bookID

	var item model.ShelfItem
	result := r.db.First(&item, "id = ?", id)

	return item, result.Error == nil
}

// Create a shelf item
func (r *Repo) Create(item model.ShelfItem) error {
	if item.UserID == "" {
		return ErrUserIDRequired
	}
	if item.BookID == "" {
		return ErrBookIDRequired
	}

	item.ID = item.UserID + ":" + item.BookID

	var existing model.ShelfItem
	err := r.db.First(&existing, "id = ?", item.ID).Error
	if err == nil {
		return ErrShelfExists
	}

	return r.db.Create(&item).Error
}
