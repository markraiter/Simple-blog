package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/markraiter/simple-blog/config"
	"github.com/markraiter/simple-blog/internal/model"
)

type TokenPair struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	AccessExpire  time.Time
	RefreshExpire time.Time
}

// NewTokenPair generates new custom struct TokenPair and returns it.
//
// In case of error occurs it throws an error.
func NewTokenPair(cfg config.Auth, user *model.User) (*TokenPair, error) {
	const operation = "jwt.NewTokenPair"

	accessExpire := time.Now().Add(cfg.AccessTTL)
	refreshExpire := time.Now().Add(cfg.RefreshTTL)

	accessToken, err := NewToken(cfg, *user, cfg.AccessTTL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	refreshToken, err := NewToken(cfg, *user, cfg.RefreshTTL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	tokenPair := TokenPair{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		AccessExpire:  accessExpire,
		RefreshExpire: refreshExpire,
	}

	return &tokenPair, nil
}

// NewToken generates new JWT token and returns signedString.
//
// In case of error occurs it throws an error.
func NewToken(cfg config.Auth, user model.User, duration time.Duration) (string, error) {
	const operation = "jwt.NewToken"

	token := jwt.New(jwt.SigningMethodHS256)

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", model.ErrTypeAssert
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

// NewRegisterToken generates new JWT token and returns signedString.
//
// In case of error occurs it throws an error.
func NewRegisterToken(cfg config.Auth, user *model.User, duration time.Duration) (string, error) {
	const operation = "jwt.NewRegisterToken"

	token := jwt.New(jwt.SigningMethodHS256)

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", model.ErrTypeAssert
	}

	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["password"] = user.Password
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
			return nil, model.ErrInvalidSigningMethod
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("accessToken throws an error during parsing: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", model.ErrInvalidToken
	}

	userID, ok := claims["uid"].(string)
	if !ok {
		return "", model.ErrNotFoundInTokenClaims
	}

	return userID, nil
}
