package service

import (
	"fmt"
	"log/slog"

	"github.com/markraiter/simple-blog/internal/model"
)

type UserSaver interface {
	SaveUser(user *model.User) (int, error)
}

type UserProvider interface {
	UserByEmail(email string) (*model.User, error)
}

type AuthService struct {
	log      *slog.Logger
	saver    UserSaver
	provider UserProvider
}

func (as *AuthService) RegisterUser(user *model.UserRequest) (int, error) {
	const operation = "AuthService.RegisterUser"

	log := as.log.With(slog.String("operation", operation))

	log.Info("checking if user already exists")

	_, err := as.provider.UserByEmail(user.Email)
	if err == nil {
		return 0, model.ErrUserAlreadyExists
	}

	if err != model.ErrUserNotFound {
		return 0, fmt.Errorf("error checking user: %w", err)
	}

	log.Info("this is a new user")
	log.Info("saving user")

	userResp := model.User{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}

	id, err := as.saver.SaveUser(&userResp)
	if err != nil {
		return 0, fmt.Errorf("error saving user: %w", err)
	}

	log.Info("user saved")

	return id, nil
}
