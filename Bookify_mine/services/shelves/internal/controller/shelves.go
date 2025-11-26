package controller

import (
	"errors"

	"bookify/services/shelves/internal/model"
	"bookify/services/shelves/internal/repository/memory"
)

var (
	ErrAlreadyExists  = errors.New("shelf item already exists")
	ErrUserIDRequired = errors.New("userId required")
	ErrBookIDRequired = errors.New("bookId required")
)

// NEW: Interface for both repos
type Repository interface {
	Get(user, book string) (model.ShelfItem, bool)
	Create(item model.ShelfItem) error
}

type Shelves struct {
	repo Repository
}

func NewShelves(repo Repository) *Shelves {
	return &Shelves{repo: repo}
}

func (c *Shelves) Get(user, book string) (model.ShelfItem, bool) {
	return c.repo.Get(user, book)
}

func (c *Shelves) Create(item model.ShelfItem) (model.ShelfItem, error) {
	if item.UserID == "" {
		return model.ShelfItem{}, ErrUserIDRequired
	}
	if item.BookID == "" {
		return model.ShelfItem{}, ErrBookIDRequired
	}

	if err := c.repo.Create(item); err != nil {
		if errors.Is(err, memory.ErrShelfExists) {
			return model.ShelfItem{}, ErrAlreadyExists
		}
		if errors.Is(err, memory.ErrUserIDRequired) {
			return model.ShelfItem{}, ErrUserIDRequired
		}
		if errors.Is(err, memory.ErrBookIDRequired) {
			return model.ShelfItem{}, ErrBookIDRequired
		}
		return model.ShelfItem{}, err
	}

	return item, nil
}
