package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markraiter/simple-blog/models"
)

type LoginInput struct {
	Email    string
	Password string
}

func (h *Handler) signUp(c *gin.Context) {
	user := new(models.User)

	if err := c.Bind(&user); err != nil {
		log.Printf("incorrect user data: %v\n", err)
		c.String(http.StatusBadRequest, "incorrect user data")
		return
	}

	email := h.storage.Authentication.GetEmail(user.Email)

	if user.Email == email {
		log.Printf("this user is already exists: %v", email)
		c.String(http.StatusConflict, "this user is already exists")
		return
	}

	id, err := h.storage.Authentication.Create(user)
	if err != nil {
		log.Printf("error creating user in the database: %v\n", err)
		c.String(http.StatusInternalServerError, "error creating user in the database")
		return
	}

	log.Printf("user %v successfully created by id %v\n", user.Email, id)
	c.String(http.StatusCreated, "user successfully created by id %d", id)
}

func (h *Handler) signIn(c *gin.Context) {
	var input LoginInput

	if err := c.Bind(&input); err != nil {
		log.Printf("incorrect input: %v\n", err)
		c.String(http.StatusBadRequest, "incorrect input")
		return
	}

	token, err := h.storage.Authentication.GenerateToken(input.Email, input.Password)
	if err != nil {
		log.Printf("Incorrect email or password: %v", err)
		c.String(http.StatusNotFound, "Incorrect email or password")
		return
	}

	log.Printf("Successfully got token: %s\n", token)
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
