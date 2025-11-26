package postgres

import (
	"errors"

	"bookify/services/library/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrAlreadyExists = errors.New("book already exists")
	ErrIDRequired    = errors.New("book id required")
)

type Repo struct {
	db *gorm.DB
}

func New(dsn string) (*Repo, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&model.Book{}); err != nil {
		return nil, err
	}

	return &Repo{db: db}, nil
}

func (r *Repo) Get(id string) (model.Book, bool) {
	var book model.Book
	result := r.db.First(&book, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.Book{}, false
	}

	return book, true
}

func (r *Repo) Create(book model.Book) error {
	if book.ID == "" {
		return ErrIDRequired
	}
	var existing model.Book
	if err := r.db.First(&existing, "id = ?", book.ID).Error; err == nil {
		return ErrAlreadyExists
	}

	return r.db.Create(&book).Error
}
