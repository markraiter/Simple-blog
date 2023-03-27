package main

import (
	"github.com/markraiter/simple-blog/internal/initializers"
	"github.com/markraiter/simple-blog/internal/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.Post{}, &models.Comment{}, models.User{})
}