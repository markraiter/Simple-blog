package models

type Post struct {
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
