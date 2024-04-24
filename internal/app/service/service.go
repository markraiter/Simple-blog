package service

import "log/slog"

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
