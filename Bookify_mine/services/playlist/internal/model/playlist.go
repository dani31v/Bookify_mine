package model

type Song struct {
	Title  string `json:"title" gorm:"column:title"`
	Artist string `json:"artist" gorm:"column:artist"`
}

type Playlist struct {
	ID     string `json:"id" gorm:"primaryKey"`
	BookID string `json:"bookId" gorm:"column:book_id"`
	Tracks []Song `json:"tracks" gorm:"-"`
}
