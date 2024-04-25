package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/markraiter/simple-blog/config"
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

	log.Info("hashing password")

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	userResp := model.User{
		Username: user.Username,
		Password: string(passHash),
		Email:    user.Email,
	}

	log.Info("saving user")

	id, err := as.saver.SaveUser(ctx, &userResp)
	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			return 0, fmt.Errorf("%s: %w", operation, model.ErrUserAlreadyExists)
		}

		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("user saved")

	return id, nil
}

func (as *AuthService) Login(ctx context.Context, cfg config.Auth, email, password string) (*jwt.TokenPair, error) {
	const operation = "AuthService.Login"

	log := as.log.With(slog.String("operation", operation))

	log.Info("getting user")

	user, err := as.provider.User(ctx, email)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, fmt.Errorf("%s: %w", operation, model.ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("comparing passwords")

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("%s: %w", operation, model.ErrInvalidCredentials)
	}

	tokenPair, err := jwt.NewTokenPair(cfg, user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("user logged in")

	return tokenPair, nil
}
