package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/markraiter/simple-blog/config"
	"github.com/markraiter/simple-blog/internal/model"
)

var (
	ErrTypeAssert            = errors.New("type assertion failed")
	ErrInvalidSigningMethod  = errors.New("invalid signing method")
	ErrInvalidToken          = errors.New("invalid token")
	ErrNotFoundInTokenClaims = errors.New("not found in token claims")
)

type TokenPair struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	AccessExpire  time.Time
	RefreshExpire time.Time
}

// NewToken generates new JWT token and returns signedString.
//
// In case of error occurs it throws an error.
func NewToken(cfg config.Auth, user *model.User, duration time.Duration) (string, error) {
	const operation = "jwt.NewToken"

	token := jwt.New(jwt.SigningMethodHS256)

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrTypeAssert
	}

	claims["uid"] = user.ID
	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(cfg.SigningKey))
	if err != nil {
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	return tokenString, nil
}

// ParseToken parses the JWT token and returns the user ID.
//
// If the token is invalid, returns an error.
// If the token is valid, returns the user ID.
func ParseToken(tokenString, signingKey string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("accessToken throws an error during parsing: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", ErrInvalidToken
	}

	userID, ok := claims["uid"].(string)
	if !ok {
		return "", ErrNotFoundInTokenClaims
	}

	return userID, nil
}
