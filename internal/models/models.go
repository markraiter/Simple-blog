package models

import "gorm.io/gorm"

type Models struct {
	Post    Post
	Comment Comment
	User User
}

type Post struct {
	gorm.Model
	UserID int    `json:"userId" xml:"userId"`
	Title  string `json:"title" xml:"title"`
	Body   string `json:"body" xml:"body"`
}

type Comment struct {
	gorm.Model
	PostID int    `json:"postId" xml:"postId"`
	Name   string `json:"name" xml:"name"`
	Email  string `json:"email" xml:"email"`
	Body   string `json:"body" xml:"body"`
}

type User struct {
	gorm.Model
	Email    string `json:"email" xml:"email"`
	Password string `json:"password"`
}
