package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markraiter/simple-blog/models"
)

func (h *Handler) getAllPosts(c echo.Context) error {
	return nil
}

func (h *Handler) getPostByID(c echo.Context) error {
	return nil
}

func (h *Handler) createPost(c echo.Context) error {
	userID, err := getUserId(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	var input models.Post

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	id, err := h.services.Posts.Create(userID, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) updatePost(c echo.Context) error {
	return nil
}

func (h *Handler) deletePost(c echo.Context) error {
	return nil
}
