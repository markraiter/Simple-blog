package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/markraiter/simple-blog/config"
	_ "github.com/markraiter/simple-blog/docs"
	"github.com/markraiter/simple-blog/internal/app/service"
	"github.com/markraiter/simple-blog/internal/lib/sl"
	"github.com/markraiter/simple-blog/internal/model"
)

type Auth interface {
	RegisterUser(ctx context.Context, user *model.UserRequest) (int, error)
	Login(ctx context.Context, cfg config.Auth, email, password string) (string, error)
}

type AuthHandler struct {
	log      *slog.Logger
	validate *validator.Validate
	service  Auth
}

// @Summary Register user
// @Description Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.UserRequest true "User data"
// @Success 201 {string} string "User ID"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/auth/register [post]
func (ah *AuthHandler) RegisterUser(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const operation = "AuthHandler.RegisterUser"

		log := ah.log.With(slog.String("operation", operation))

		log.Debug("parsing request")

		var userReq model.UserRequest

		if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
			log.Warn("error parsing request", sl.Err(err))
			http.Error(w, "error parsing request", http.StatusBadRequest)

			return
		}

		log.Debug("validating user")

		if err := ah.validate.Struct(userReq); err != nil {
			log.Warn("error validating user", sl.Err(err))
			http.Error(w, "error validating user", http.StatusBadRequest)

			return
		}

		log.Debug("registering user")

		id, err := ah.service.RegisterUser(ctx, &userReq)
		if err != nil {
			if errors.Is(err, service.ErrAlreadyExists) {
				log.Warn("user already exists", sl.Err(err))
				http.Error(w, "user already exists", http.StatusBadRequest)

				return
			}

			log.Error("error registering user", sl.Err(err))
			http.Error(w, "error registering user", http.StatusInternalServerError)

			return
		}

		log.Debug("user registered")

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(strconv.Itoa(id)))
	}
}

// @Summary Login
// @Description Login
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.LoginRequest true "User data"
// @Success 200 {string} string "Token"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/auth/login [post]
func (ah *AuthHandler) Login(ctx context.Context, cfg config.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "AuthHandler.Login"

		log := ah.log.With(slog.String("operation", operation))

		log.Debug("parsing request")

		var userReq model.LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
			log.Warn("error parsing request", sl.Err(err))
			http.Error(w, "error parsing request", http.StatusBadRequest)

			return
		}

		log.Debug("validating user")

		if err := ah.validate.Struct(userReq); err != nil {
			log.Warn("error validating user", sl.Err(err))
			http.Error(w, "error validating user", http.StatusBadRequest)

			return
		}

		log.Debug("logging in")

		token, err := ah.service.Login(ctx, cfg, userReq.Email, userReq.Password)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				log.Warn("user not found", sl.Err(err))
				http.Error(w, "user not found", http.StatusBadRequest)

				return
			}

			if errors.Is(err, service.ErrInvalidCredentials) {
				log.Warn("invalid credentials", sl.Err(err))
				http.Error(w, "invalid credentials", http.StatusBadRequest)

				return
			}

			log.Error("error logging in", sl.Err(err))
			http.Error(w, "error logging in", http.StatusInternalServerError)

			return
		}

		log.Debug("logged in")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(token)
	}
}
