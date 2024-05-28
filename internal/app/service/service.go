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

type Service struct {
	AuthService
}

func New(log *slog.Logger, auth AuthStorage) *Service {
	return &Service{
		AuthService{
			log:      log,
			saver:    auth,
			provider: auth,
		},
	}
}
