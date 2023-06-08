package models

type Post struct {
	UserID int    `json:"userId" db:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
