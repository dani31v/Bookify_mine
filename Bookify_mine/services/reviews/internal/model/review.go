
package model

type Review struct {
	ID     string `json:"id"`
	BookID string `json:"bookId"`
	UserID string `json:"userId"`
	Rating int    `json:"rating"`
	Text   string `json:"text"`
}
