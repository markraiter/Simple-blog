package models

import "errors"

type Comment struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id" db:"user_id"`
	PostID uint   `json:"post_id" db:"post_id"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

type UpdateCommentInput struct {
	Email *string `json:"email"`
	Body  *string `json:"body"`
}

func (i *UpdateCommentInput) Validate() error {
	if i.Email == nil && i.Body == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
