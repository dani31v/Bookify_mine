package memory

import (
	"errors"
	"sync"

	"bookify/services/library/internal/model"
)

var (
	ErrBookExists     = errors.New("book already exists")
	ErrBookIDRequired = errors.New("book id required")
)

type Repo struct {
	mu   sync.RWMutex
	data map[string]model.Book
}

func New() *Repo {
	return &Repo{data: map[string]model.Book{
		"book-1": {ID: "book-1", Title: "Pride and Prejudice", Author: "Jane Austen", Pages: 364, Edition: "2nd"},
	}}
}

func (r *Repo) Get(id string) (model.Book, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.data[id]
	return b, ok
}

func (r *Repo) Create(book model.Book) error {
	if book.ID == "" {
		return ErrBookIDRequired
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[book.ID]; exists {
		return ErrBookExists
	}

	r.data[book.ID] = book
	return nil
}
