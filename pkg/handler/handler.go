package handler

import (
	"github.com/gin-gonic/gin"
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

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
	}

	api := router.Group("/api", h.userIdentity)
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

	return router
}
