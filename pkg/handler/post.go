package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/markraiter/simple-blog/models"
)

// Swagger comment for getAllPosts
func (h *Handler) getAllPosts(c *gin.Context) {
	posts, err := h.services.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// Swagger comment for getPostByID
func (h *Handler) getPostByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id param",
		})
		return
	}

	post, err := h.services.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, post)
}

// Swagger comment for createPost
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

// Swagger comment for updatePost
func (h *Handler) updatePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id param",
		})
		return
	}

	var input models.UpdatePostInput

	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := h.services.Update(id, input); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, "Updated!")
}

// Swagger comment for deletePost
func (h *Handler) deletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id param",
		})
		return
	}

	if err := h.services.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, "Deleted!")
}
