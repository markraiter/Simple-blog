package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/markraiter/simple-blog/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserSaver interface {
	SaveUser(ctx context.Context, user *model.User) (int, error)
}

type UserProvider interface {
	UserByEmail(ctx context.Context, email string) (*model.User, error)
}

type AuthService struct {
	log      *slog.Logger
	saver    UserSaver
	provider UserProvider
}

func (as *AuthService) RegisterUser(ctx context.Context, user *model.UserRequest) (int, error) {
	const operation = "AuthService.RegisterUser"

	log := as.log.With(slog.String("operation", operation))

	log.Info("checking if user already exists")

	_, err := as.provider.UserByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			return 0, fmt.Errorf("%s: %w", operation, model.ErrUserAlreadyExists)
		}
	}

	log.Info("this is a new user")
	log.Info("saving user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	userResp := model.User{
		Username: user.Username,
		Password: string(passHash),
		Email:    user.Email,
	}

	id, err := as.saver.SaveUser(ctx, &userResp)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("user saved")

	return id, nil
}
