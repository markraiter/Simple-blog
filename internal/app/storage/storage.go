package storage

import "errors"

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
	ErrNotAllowed    = errors.New("not allowed")
	ErrPostNotExists = errors.New("post with such ID does not exist")
)
