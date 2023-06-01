package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/markraiter/simple-blog/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *echo.Echo {
	e := echo.New()

	auth := e.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
	}

	api := e.Group("/api", h.userIdentity)
	{
		posts := api.Group("/posts")
		{
			posts.GET("/", h.getAllPosts)
			posts.GET("/:id", h.getPostByID)
			posts.POST("/", h.createPost)
			posts.PUT("/:id", h.updatePost)
			posts.DELETE("/:id", h.deletePost)
		}

		comments := api.Group("/comments")
		{
			comments.GET("/", h.getAllComments)
			comments.GET("/:id", h.getCommentByID)
			comments.POST("/", h.createComment)
			comments.PUT("/:id", h.updateComment)
			comments.DELETE("/:id", h.deleteComment)
		}
	}

	return e
}
