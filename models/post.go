package models

type Post struct {
	UserID int    `json:"userId" xml:"userId"`
	Title  string `json:"title" xml:"title"`
	Body   string `json:"body" xml:"body"`
}
