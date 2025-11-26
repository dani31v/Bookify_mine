package postgres

import (
	"errors"
	"fmt"

	"bookify/services/reviews/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrReviewExists     = errors.New("review already exists")
	ErrReviewIDRequired = errors.New("id required")
	ErrBookIDRequired   = errors.New("bookId required")
)

type Repo struct {
	db *gorm.DB
}

func New(dsn string) (*Repo, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgres connect error: %w", err)
	}

	if err := db.AutoMigrate(&model.Review{}); err != nil {
		return nil, fmt.Errorf("migrate error: %w", err)
	}

	return &Repo{db: db}, nil
}

func (r *Repo) ByBook(bookID string) []model.Review {
	var reviews []model.Review
	r.db.Where("book_id = ?", bookID).Find(&reviews)
	return reviews
}

func (r *Repo) Create(review model.Review) error {
	if review.ID == "" {
		return ErrReviewIDRequired
	}
	if review.BookID == "" {
		return ErrBookIDRequired
	}

	var existing model.Review
	err := r.db.First(&existing, "id = ?", review.ID).Error
	if err == nil {
		return ErrReviewExists
	}

	return r.db.Create(&review).Error
}
