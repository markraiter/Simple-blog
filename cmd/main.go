package main

import (
	"github.com/markraiter/simple-blog/cmd/app"
)

// @title Simple-Blog API
// @version 1.0
// @description API Server for Simple Blog Application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	app.Start()
}
