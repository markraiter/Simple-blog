package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			return c.String(http.StatusUnauthorized, "authorization header missing")
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v/n", t.Header["alg"])
			}

			return []byte("secret"), nil
		})

		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string {
				"error": err.Error(),
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := uint(claims["id"].(float64))
			c.Set("userID", userID)

			return next(c)
		} else {
			return c.JSON(http.StatusUnauthorized, map[string]string {
				"error": "invalid token",
			})
		}
	}
}