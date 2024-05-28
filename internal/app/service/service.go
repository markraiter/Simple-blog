package service

import (
	"errors"
	"log/slog"
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
	l *slog.Logger,
	a AuthStorage,
	p PostStorage,
) *Service {
	return &Service{
		AuthService{
			log:      l,
			saver:    a,
			provider: a,
		},
		PostService{
			log:       l,
			saver:     p,
			provider:  p,
			processor: p,
		},
	}
}
