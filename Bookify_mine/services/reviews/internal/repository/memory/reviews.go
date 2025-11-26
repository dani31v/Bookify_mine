package memory

import (
	"errors"
	"sync"

	"bookify/services/reviews/internal/model"
)

var (
	ErrReviewExists     = errors.New("review already exists")
	ErrReviewIDRequired = errors.New("id required")
	ErrBookIDRequired   = errors.New("bookId required")
)

type Repo struct {
	mu   sync.RWMutex
	data []model.Review
	ids  map[string]struct{}
}

func New() *Repo {
	data := []model.Review{
		{ID: "rev-1", BookID: "book-1", UserID: "dani", Rating: 5, Text: "Hermoso."},
		{ID: "rev-2", BookID: "book-1", UserID: "santi", Rating: 4, Text: "Clásico."},
	}
	ids := make(map[string]struct{}, len(data))
	for _, review := range data {
		ids[review.ID] = struct{}{}
	}
	return &Repo{data: data, ids: ids}
}

func (r *Repo) ByBook(id string) []model.Review {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]model.Review, 0, len(r.data))
	for _, v := range r.data {
		if v.BookID == id {
			res = append(res, v)
		}
	}
	return res
}

func (r *Repo) Create(review model.Review) error {
	if review.ID == "" {
		return ErrReviewIDRequired
	}
	if review.BookID == "" {
		return ErrBookIDRequired
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.ids[review.ID]; exists {
		return ErrReviewExists
	}

	r.data = append(r.data, review)
	r.ids[review.ID] = struct{}{}
	return nil
}
