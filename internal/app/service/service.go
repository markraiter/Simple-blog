package service

import (
	"errors"
)

var (
	ErrAlreadyExists      = errors.New("already exists")
	ErrNotFound           = errors.New("not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
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

type Service struct {
	AuthService
	PostService
}

func New(
	a AuthStorage,
	p PostStorage,
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
	}
}
