package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/markraiter/simple-blog/config"
	"github.com/markraiter/simple-blog/internal/app/storage"
	"github.com/markraiter/simple-blog/internal/lib/jwt"
	"github.com/markraiter/simple-blog/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserSaver interface {
	SaveUser(ctx context.Context, user *model.User) (int, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (*model.User, error)
}

type AuthService struct {
	log      *slog.Logger
	saver    UserSaver
	provider UserProvider
}

func (as *AuthService) RegisterUser(ctx context.Context, user *model.UserRequest) (int, error) {
	const operation = "AuthService.RegisterUser"

	log := as.log.With(slog.String("operation", operation))

	log.Debug("hashing password")

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	userResp := model.User{
		Username: user.Username,
		Password: string(passHash),
		Email:    user.Email,
	}

	log.Debug("saving user")

	id, err := as.saver.SaveUser(ctx, &userResp)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return 0, fmt.Errorf("%s: %w", operation, ErrAlreadyExists)
		}

		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	log.Debug("user saved")

	return id, nil
}

func (as *AuthService) Login(ctx context.Context, cfg config.Auth, email, password string) (string, error) {
	const operation = "AuthService.Login"

	log := as.log.With(slog.String("operation", operation))

	log.Debug("getting user")

	user, err := as.provider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return "", fmt.Errorf("%s: %w", operation, ErrNotFound)
		}

		return "", fmt.Errorf("%s: %w", operation, err)
	}

	log.Debug("comparing passwords")

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("%s: %w", operation, ErrInvalidCredentials)
	}

	token, err := jwt.NewToken(cfg, user, cfg.AccessTTL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	log.Debug("user logged in")

	return token, nil
}
