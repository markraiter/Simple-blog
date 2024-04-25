package model

import (
	"errors"
	"log/slog"
)

var (
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrNotFoundInTokenClaims = errors.New("not found in token claims")
	ErrInvalidToken          = errors.New("invalid token")
	ErrInvalidSigningMethod  = errors.New("invalid signing method")
	ErrTypeAssert            = errors.New("type assert error")
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
