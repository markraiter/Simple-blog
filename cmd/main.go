package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator"
	"github.com/markraiter/simple-blog/config"
	"github.com/markraiter/simple-blog/internal/app/api"
	"github.com/markraiter/simple-blog/internal/app/api/handler"
	"github.com/markraiter/simple-blog/internal/app/api/middleware"
	"github.com/markraiter/simple-blog/internal/app/service"
	"github.com/markraiter/simple-blog/internal/app/storage"
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

	db := storage.New(cfg.Postgres)

	service := service.New(log, db)

	handler := handler.New(log, validate, service)

	server := new(api.Server)

	go func() {
		if err := server.Run(cfg, handler.Router()); err != nil {
			log.Error("error occured while running http server: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	log.Info("shutting down application...")

	if err := server.Shutdown(context.TODO()); err != nil {
		log.Error("error occured while shutting down the server: " + err.Error())
	}
}
