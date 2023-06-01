package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	authHeader = "Authorization"
	userCtx    = "userID"
)

func (h *Handler) userIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get(authHeader)
		if header == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "empty auth header",
			})
		}

		tokenString := strings.Replace(header, "Bearer ", "", 1)

		userID, err := h.services.ParseToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
		}

		c.Set(userCtx, userID)

		return nil
	}
}

func getUserId(c echo.Context) (int, error) {
	id := c.Get(userCtx)

	idInt, ok := id.(int)
	if !ok {
		return 0, fmt.Errorf("user id is of invalid type")
	}

	return idInt, nil
}
