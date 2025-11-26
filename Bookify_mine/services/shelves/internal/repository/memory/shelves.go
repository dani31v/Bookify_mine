package memory

import (
	"errors"
	"sync"

	"bookify/services/shelves/internal/model"
)

var (
	ErrShelfExists    = errors.New("shelf item already exists")
	ErrUserIDRequired = errors.New("userId required")
	ErrBookIDRequired = errors.New("bookId required")
)

type Repo struct {
	mu   sync.RWMutex
	data map[string]model.ShelfItem
}

func New() *Repo {
	return &Repo{data: map[string]model.ShelfItem{
		"dani:book-1": {UserID: "dani", BookID: "book-1", Shelf: model.ToRead},
	}}
}

func (r *Repo) Get(userId, bookId string) (model.ShelfItem, bool) {
	key := userId + ":" + bookId
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.data[key]
	return v, ok
}

func (r *Repo) Create(item model.ShelfItem) error {
	if item.UserID == "" {
		return ErrUserIDRequired
	}
	if item.BookID == "" {
		return ErrBookIDRequired
	}

	key := item.UserID + ":" + item.BookID

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[key]; exists {
		return ErrShelfExists
	}

	r.data[key] = item
	return nil
}
