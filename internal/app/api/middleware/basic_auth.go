package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/markraiter/simple-blog/config"
	"github.com/markraiter/simple-blog/internal/lib/jwt"
	"github.com/markraiter/simple-blog/internal/lib/sl"
)

type contextKey string

const (
	UIDKey           contextKey = "uid"
	RefreshStringKey contextKey = "refreshString"
	EmailKey         contextKey = "email"
	UsernameKey      contextKey = "username"
)

func BasicAuth(cfg config.Auth, log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const operation = "middleware.BasicAuth"

			l := log.With(slog.String("operation", operation))

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				l.Warn("Authorization header is required")
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}
			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

			tokenClaims, err := jwt.ParseToken(tokenString, cfg.SigningKey)
			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					l.Warn("token expired", sl.Err(err))
					http.Error(w, "Token expired", http.StatusUnauthorized)
					return
				}

				l.Warn("error parsing token", sl.Err(err))
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UIDKey, tokenClaims.UID)
			ctx = context.WithValue(ctx, RefreshStringKey, tokenString)
			ctx = context.WithValue(ctx, EmailKey, tokenClaims.Email)
			ctx = context.WithValue(ctx, UsernameKey, tokenClaims.Username)

			// spew.Dump(ctx)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
