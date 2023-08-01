package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/markraiter/simple-blog/internal/storage"
)

type Handler struct {
	storage *storage.Storage
}

func NewHandler(storage *storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.JWTMiddleware)
	{

		posts := api.Group("/posts")
		{
			posts.GET("/all", h.getAllPosts)
			posts.GET("", h.filterPostsByUser)
			posts.GET("/:id", h.getPostByID)
			posts.POST("", h.createPost)
			posts.PATCH("/:id", h.updatePost)
			posts.DELETE("/:id", h.deletePost)
		}

		comments := api.Group("/comments")
		{
			comments.GET("/all", h.getAllComments)
			comments.GET("/post", h.filterCommentsByPost)
			comments.GET("/user", h.filterCommentsByUser)
			comments.GET("/:id", h.getCommentByID)
			comments.POST("", h.createComment)
			comments.PATCH("/:id", h.updateComment)
			comments.DELETE("/:id", h.deleteComment)
		}
	}

	return router
}
