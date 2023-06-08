package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markraiter/simple-blog/models"
)

// swagger comments for register...
func (h *Handler) register(c *gin.Context) {
	var input models.User

	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
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

// swagger comments for login
func (h *Handler) login(c *gin.Context) {
	var input models.User

	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
