package app

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/markraiter/simple-blog/internal/initializers"
	"github.com/markraiter/simple-blog/pkg/handlers"
)

func Start() {
	e := echo.New()

	e.GET("/", handlers.Hello)

	e.POST("/registration", auth.Register(initializers.DB))
	e.POST("/login", auth.Login)

	e.GET("/posts", handlers.GetPosts(initializers.DB))
	e.GET("/posts:id", handlers.GetPost(initializers.DB))
	e.POST("/posts", handlers.CreatePost(initializers.DB))
	e.PUT("/posts:id", handlers.UpdatePost(initializers.DB))
	e.DELETE("/posts:id", handlers.DeletePost(initializers.DB))

	e.GET("/comments", handlers.GetComments(initializers.DB))
	e.GET("/comments:id", handlers.GetComment(initializers.DB))
	e.POST("/comments", handlers.CreateComment(initializers.DB))
	e.PUT("/comments:id", handlers.UpdateComment(initializers.DB))
	e.DELETE("/comments:id", handlers.DeleteComment(initializers.DB))

	log.Fatal(e.Start(":8080"))

}