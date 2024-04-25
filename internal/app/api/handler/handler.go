package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
	httpSwagger "github.com/swaggo/http-swagger"
)

type AuthService interface {
	UserSaver
	UserProvider
}

type Handler struct {
	Healthcheck
	AuthHandler
}

func New(log *slog.Logger, validate *validator.Validate, auth AuthService) *Handler {
	return &Handler{
		Healthcheck{log: log},
		AuthHandler{log: log, validate: validate, saver: auth, provider: auth},
	}
}

func (h *Handler) Router(ctx context.Context) http.Handler {
	m := http.NewServeMux()

	m.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	m.Handle("GET /health", h.APIHealth())
	m.Handle("POST /auth/register", h.RegisterUser(ctx))

	return m
}
