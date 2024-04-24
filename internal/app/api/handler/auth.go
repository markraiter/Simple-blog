package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/markraiter/simple-blog/internal/model"
)

type UserSaver interface {
	RegisterUser(user *model.UserRequest) (int, error)
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

func (ah *AuthHandler) RegisterUser() http.HandlerFunc {
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

		id, err := ah.saver.RegisterUser(&userReq)
		if err != nil {
			log.Error("error registering user", model.Err(err))
			http.Error(w, "error registering user", http.StatusInternalServerError)
			return
		}

		log.Info("user registered")

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(strconv.Itoa(id)))
	}
}
