package controller

import (
	"errors"

	"bookify/services/reviews/internal/model"
	"bookify/services/reviews/internal/repository/memory"
)

var (
	ErrAlreadyExists  = errors.New("review already exists")
	ErrIDRequired     = errors.New("id required")
	ErrBookIDRequired = errors.New("bookId required")
)

type Repository interface {
	ByBook(id string) []model.Review
	Create(review model.Review) error
}

type Reviews struct {
	repo Repository
}

func NewReviews(repo Repository) *Reviews {
	return &Reviews{repo: repo}
}

func (c *Reviews) ByBook(id string) []model.Review {
	return c.repo.ByBook(id)
}

func (c *Reviews) Create(review model.Review) (model.Review, error) {
	if review.ID == "" {
		return model.Review{}, ErrIDRequired
	}
	if review.BookID == "" {
		return model.Review{}, ErrBookIDRequired
	}

	if err := c.repo.Create(review); err != nil {
		if errors.Is(err, memory.ErrReviewExists) {
			return model.Review{}, ErrAlreadyExists
		}
		if errors.Is(err, memory.ErrReviewIDRequired) {
			return model.Review{}, ErrIDRequired
		}
		if errors.Is(err, memory.ErrBookIDRequired) {
			return model.Review{}, ErrBookIDRequired
		}
		return model.Review{}, err
	}

	return review, nil
}
