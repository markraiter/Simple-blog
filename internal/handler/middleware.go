package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const authHeader = "Authorization"

// JWTMiddleware is authentication middleware for checking the correct token
func (h *Handler) JWTMiddleware(c *gin.Context) {
	authHeader := c.GetHeader(authHeader)
	if authHeader == "" {
		log.Printf("authorization header missing")
		c.String(http.StatusUnauthorized, "authorization header missing")
		c.Abort()
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method: %v", t.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte("secret"), nil
	})

	if err != nil {
		log.Printf("error parsing token: %v", err)
		c.String(http.StatusUnauthorized, "error parsing token")
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["id"]
		c.Set("user_id", userID)

		c.Next()
	} else {
		log.Printf("invalid token")
		c.String(http.StatusUnauthorized, "invalid token")
	}
}
