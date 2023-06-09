package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/markraiter/simple-blog/models"
)

func (h *Handler) getAllComments(c *gin.Context) {
	comments, err := h.services.Comments.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, comments)
}

func (h *Handler) getCommentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id param",
		})
		return
	}

	comment, err := h.services.Comments.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *Handler) createComment(c *gin.Context) {
	postID, err := getPostId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	var input models.Comment

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	id, err := h.services.Comments.Create(postID, input)
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

func (h *Handler) updateComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id param",
		})
		return
	}

	var input models.UpdateCommentInput

	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := h.services.Comments.Update(id, input); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"message": "Updated!",
	})
}

func (h *Handler) deleteComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id param",
		})
		return
	}

	if err := h.services.Comments.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, map[string]string{
		"message": "Deleted!",
	})
}
