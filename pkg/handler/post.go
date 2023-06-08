package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markraiter/simple-blog/models"
)

func (h *Handler) getAllPosts(c *gin.Context) {

}

func (h *Handler) getPostByID(c *gin.Context) {

}

func (h *Handler) createPost(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	var input models.Post

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	id, err := h.services.Posts.Create(userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) updatePost(c *gin.Context) {

}

func (h *Handler) deletePost(c *gin.Context) {

}
