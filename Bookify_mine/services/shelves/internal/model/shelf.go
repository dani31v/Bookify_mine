package model

type Shelf string

const (
	ToRead  Shelf = "to-read"
	Reading Shelf = "reading"
	Done    Shelf = "finished"
)

type ShelfItem struct {
	ID     string `gorm:"primaryKey" json:"-"` // userId + ":" + bookId
	UserID string `json:"userId"`
	BookID string `json:"bookId"`
	Shelf  Shelf  `json:"shelf"`
}
