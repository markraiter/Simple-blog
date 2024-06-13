package service

import (
	"errors"
)

var (
	ErrAlreadyExists      = errors.New("already exists")
	ErrNotFound           = errors.New("not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotAllowed         = errors.New("user is not allowed to perform this operation")
	ErrPostNotExists      = errors.New("post with such ID does not exist")
)

type AuthStorage interface {
	UserSaver
	UserProvider
}

type PostStorage interface {
	PostSaver
	PostProvider
	PostProcessor
}

type CommentStorage interface {
	CommentSaver
	CommentProvider
	CommentProcessor
}

type Service struct {
	AuthService
	PostService
	CommentService
}

func New(
	a AuthStorage,
	p PostStorage,
	c CommentStorage,
) *Service {
	return &Service{
		AuthService{
			saver:    a,
			provider: a,
		},
		PostService{
			saver:     p,
			provider:  p,
			processor: p,
		},
		CommentService{
			saver:     c,
			provider:  c,
			processor: c,
		},
	}
}
