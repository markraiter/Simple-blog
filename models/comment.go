package models

type Comment struct {
	PostID int    `json:"postId" xml:"postId"`
	Name   string `json:"name" xml:"name"`
	Email  string `json:"email" xml:"email"`
	Body   string `json:"body" xml:"body"`
}
