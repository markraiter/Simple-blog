package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	_ "github.com/markraiter/simple-blog/docs"
	"github.com/markraiter/simple-blog/internal/model"
)

type UserSaver interface {
	RegisterUser(ctx context.Context, user *model.UserRequest) (int, error)
}

type UserProvider interface {
	// UserByEmail(email string) (*model.User, error)
}

type AuthHandler struct {
	log      *slog.Logger
	validate *validator.Validate
	saver    UserSaver
	provider UserProvider
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
// @Router /auth/register [post]
func (ah *AuthHandler) RegisterUser(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const operation = "AuthHandler.RegisterUser"

		log := ah.log.With(slog.String("operation", operation))

		log.Info("parsing request")

		var userReq model.UserRequest

		if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
			log.Error("error parsing request", model.Err(err))
			http.Error(w, "error parsing request", http.StatusBadRequest)

			return
		}

		log.Info("validating user")

		if err := ah.validate.Struct(userReq); err != nil {
			log.Error("error validating user", model.Err(err))
			http.Error(w, "error validating user", http.StatusBadRequest)

			return
		}

		log.Info("registering user")

		id, err := ah.saver.RegisterUser(ctx, &userReq)
		if err != nil {
			if errors.Is(err, model.ErrUserAlreadyExists) {
				log.Error("user already exists", model.Err(err))
				http.Error(w, "user already exists", http.StatusBadRequest)

				return
			}

			log.Error("error registering user", model.Err(err))
			http.Error(w, "error registering user", http.StatusInternalServerError)

			return
		}

		log.Info("user registered")

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(strconv.Itoa(id)))
	}
}
