package models

import "errors"

type Comment struct {
	PostID int    `json:"postId" xml:"postId"`
	Name   string `json:"name" xml:"name"`
	Email  string `json:"email" xml:"email"`
	Body   string `json:"body" xml:"body"`
}

type UpdateCommentInput struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Body  *string `json:"body"`
}

func (i UpdateCommentInput) Validate() error {
	if i.Name == nil && i.Email == nil && i.Body == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
