package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/markraiter/simple-blog/config"
	httpSwagger "github.com/swaggo/http-swagger"
)

type AuthService interface {
	Auth
}

type PostService interface {
	PostSaver
	PostProvider
	PostProcessor
}

type Handler struct {
	Healthcheck
	AuthHandler
	PostHandler
}

// The response struct is used to send a message back to the client.
type response struct {
	Message string `json:"message"`
}

func New(
	l *slog.Logger,
	v *validator.Validate,
	a AuthService,
	p PostService,
) *Handler {
	return &Handler{
		Healthcheck{log: l},
		AuthHandler{
			log:      l,
			validate: v,
			service:  a,
		},
		PostHandler{
			log:       l,
			validate:  v,
			saver:     p,
			provider:  p,
			processor: p,
		},
	}
}

func (h *Handler) Router(ctx context.Context, cfg config.Config) http.Handler {
	m := http.NewServeMux()

	m.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	m.Handle("GET /health", h.APIHealth())
	{
		m.Handle("POST /api/auth/register", h.RegisterUser(ctx))
		m.Handle("POST /api/auth/login", h.Login(ctx, cfg.Auth))
	}

	{
		m.Handle("POST /api/posts", h.CreatePost(ctx))
		// m.Handle("GET api/posts", h.GetPosts(ctx))
		// m.Handle("GET api/posts/{id}", h.GetPost(ctx))
		// m.Handle("PUT api/posts/{id}", h.UpdatePost(ctx))
		// m.Handle("DELETE api/posts/{id}", h.DeletePost(ctx))
	}

	return m
}
