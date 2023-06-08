package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authHeader = "Authorization"
	userCtx    = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authHeader)
	if header == "" {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "empty auth header",
		})
		return
	}

	tokenString := strings.Replace(header, "Bearer ", "", 1)

	userID, err := h.services.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.Set(userCtx, userID)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
