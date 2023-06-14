package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/markraiter/simple-blog/models"
)

// @Summary Get All Posts
// @Security ApiKeyAuth
// @Tags posts
// @Description get all posts
// @ID get-all-posts
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.Posts
// @Failure 400,404 {object} error
// @Failure 500 {object} error
// @Failure default {object} error
// @Router /api/posts [get]
func (h *Handler) getAllPosts(c *gin.Context) {
	posts, err := h.services.Posts.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// @Summary Get Post By Id
// @Security ApiKeyAuth
// @Tags posts
// @Description get post by id
// @ID get-post-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Post
// @Failure 400,404 {object} error
// @Failure 500 {object} error
// @Failure default {object} error
// @Router /api/posts/:id [get]
func (h *Handler) getPostByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id param",
		})
		return
	}

	post, err := h.services.Posts.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, post)
}

// @Summary Create post
// @Security ApiKeyAuth
// @Tags posts
// @Description create post
// @ID create-post
// @Accept  json
// @Produce  json
// @Param input body models.Post true "post info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/posts [post]
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

	if err := h.services.Posts.Update(id, input); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"message": "Updated!",
	})
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

	if err := h.services.Posts.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, map[string]string{
		"message": "Deleted!",
	})
}
