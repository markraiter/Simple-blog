package app

import (
	"github.com/labstack/echo/v4"
	"github.com/markraiter/simple-blog/internal/initializers"
	"github.com/markraiter/simple-blog/pkg/auth"
	"github.com/markraiter/simple-blog/pkg/handlers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func Start() {
	e := echo.New()

	// Groups
	// authGroup := e.Group("/api")
	// authGroup.Use(middlewares.JWTMiddleware)
	// postGroup := authGroup.Group("/v1/posts")
	// commentsGroup := authGroup.Group("/v1/comments")

	// Test handler
	e.GET("/", handlers.Hello)

	// Registration
	e.POST("/registration", auth.Register(initializers.DB))

	// Authentification
	e.POST("/login", auth.Login(initializers.DB))

	////////////////////ENDPOINTS////////////////////

	// Operations with posts
	// postGroup.GET("", handlers.GetPosts(initializers.DB))
	// postGroup.GET(":id", handlers.GetPost(initializers.DB))
	// postGroup.POST("", handlers.CreatePost(initializers.DB))
	// postGroup.PUT(":id", handlers.UpdatePost(initializers.DB))
	// postGroup.DELETE(":id", handlers.DeletePost(initializers.DB))

	// Operations with comments
	// commentsGroup.GET("", handlers.GetComments(initializers.DB))
	// commentsGroup.GET(":id", handlers.GetComment(initializers.DB))
	// commentsGroup.POST("", handlers.CreateComment(initializers.DB))
	// commentsGroup.PUT(":id", handlers.UpdateComment(initializers.DB))
	// commentsGroup.DELETE(":id", handlers.DeleteComment(initializers.DB))

	e.Logger.Fatal(e.Start(":8080"))

}
