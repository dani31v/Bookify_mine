package controller

import (
	"errors"

	"bookify/services/library/internal/model"
	"bookify/services/library/internal/repository/memory"
)

var (
	ErrAlreadyExists = errors.New("book already exists")
	ErrIDRequired    = errors.New("book id required")
)

type Repository interface {
	Get(id string) (model.Book, bool)
	Create(book model.Book) error
}

type Books struct {
	repo Repository
}

func NewBooks(repo Repository) *Books {
	return &Books{repo: repo}
}

func (c *Books) Get(id string) (any, bool) {
	return c.repo.Get(id)
}

func (c *Books) Create(book model.Book) (model.Book, error) {
	if book.ID == "" {
		return model.Book{}, ErrIDRequired
	}

	if err := c.repo.Create(book); err != nil {
		if errors.Is(err, memory.ErrBookExists) {
			return model.Book{}, ErrAlreadyExists
		}
		if errors.Is(err, memory.ErrBookIDRequired) {
			return model.Book{}, ErrIDRequired
		}
		return model.Book{}, err
	}

	return book, nil
}
