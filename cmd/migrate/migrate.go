package migrate

import (
	"github.com/markraiter/simple-blog/internal/initializers"
	"github.com/markraiter/simple-blog/internal/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func Migrate() {
	initializers.DB.AutoMigrate(&models.Post{}, &models.Comment{}, models.User{})
	WriteToDB(initializers.DB)
}