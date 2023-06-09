package models

import "errors"

type Post struct {
	ID     int    `json:"id" db:"id"`
	UserID int    `json:"user_id" db:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type UpdatePostInput struct {
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

func (i UpdatePostInput) Validate() error {
	if i.Title == nil && i.Body == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
