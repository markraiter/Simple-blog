package app

import (
	"github.com/labstack/echo/v4"
	"github.com/markraiter/simple-blog/cmd/migrate"
	"github.com/markraiter/simple-blog/internal/initializers"
	"github.com/markraiter/simple-blog/pkg/auth"
	"github.com/markraiter/simple-blog/pkg/handlers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	migrate.Migrate()
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
	e.GET("/posts", handlers.GetPosts(initializers.DB))
	e.GET("/posts/:id", handlers.GetPostByID(initializers.DB))
	// e.POST("/posts", handlers.CreatePost(initializers.DB))
	// e.PUT("/posts/:id", handlers.UpdatePost(initializers.DB))
	// e.DELETE("/posts/:id", handlers.DeletePost(initializers.DB))

	// Operations with comments
	// e.GET("/comments", handlers.GetComments(initializers.DB))
	// e.GET("/comments/:id", handlers.GetCommentByID(initializers.DB))
	// e.POST("/comments", handlers.CreateComment(initializers.DB))
	// e.PUT("/comments/:id", handlers.UpdateComment(initializers.DB))
	// e.DELETE("/cpmments/:id", handlers.DeleteComment(initializers.DB))

	e.Logger.Fatal(e.Start(":8080"))

}
