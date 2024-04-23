package main

import (
	"github.com/go-playground/validator"
	"github.com/markraiter/simple-blog/config"
	"github.com/markraiter/simple-blog/internal/app/api/middleware"
	"github.com/markraiter/simple-blog/internal/model"
)

func main() {
	cfg := config.MustLoad()

	log := middleware.SetupLogger(cfg.Env)

	validate := validator.New()
	validate.RegisterValidation("number", model.ValidateContainsNumber, false)
	validate.RegisterValidation("upper", model.ValidateContainsUpper, false)
	validate.RegisterValidation("lower", model.ValidateContainsLower, false)
	validate.RegisterValidation("special", model.ValidateContainsSpecial, false)

	log.Info("starting application...")
	log.Info("port: " + cfg.Server.Port)

}
