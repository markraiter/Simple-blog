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
	ErrTokenExpired          = errors.New("token expired")
	ErrNotFoundInTokenClaims = errors.New("not found in token claims")
)

type TokenClaims struct {
	UID      string
	Username string
	Email    string
	Exp      int64
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
func ParseToken(tokenString, signingKey string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}

		return nil, fmt.Errorf("accessToken throws an error during parsing: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	var userID string
	if uid, ok := claims["uid"]; ok {
		switch v := uid.(type) {
		case string:
			userID = v
		case float64:
			userID = fmt.Sprintf("%.0f", v)
		}
	} else {
		return nil, ErrNotFoundInTokenClaims
	}

	var exp int64
	if expClaim, ok := claims["exp"]; ok {
		switch v := expClaim.(type) {
		case float64:
			exp = int64(v)
		case int64:
			exp = v
		default:
			return nil, ErrNotFoundInTokenClaims
		}
	} else {
		return nil, ErrNotFoundInTokenClaims
	}

	var email string
	if emailClaim, ok := claims["email"]; ok {
		switch v := emailClaim.(type) {
		case string:
			email = v
		}
	} else {
		return nil, ErrNotFoundInTokenClaims
	}

	var username string
	if usernameClaim, ok := claims["username"]; ok {
		switch v := usernameClaim.(type) {
		case string:
			username = v
		}
	} else {
		return nil, ErrNotFoundInTokenClaims
	}

	tc := TokenClaims{
		UID:      userID,
		Username: username,
		Email:    email,
		Exp:      exp,
	}

	return &tc, nil
}
