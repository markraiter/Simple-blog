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

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body LoginInput true "account info"
// @Success 201 {integer} integer "User successfully registered. Returns the id of the created user."
// @Failure 400 {string} string "Invalid request or missing required fields."
// @Failure 409 {string} string "User already exists with the provided email or username."
// @Failure 500 {string} string "An unexpected error occurred on the server."
// @Router /auth/sign-up [post]
// signUp is a handler for registration. It returns the id of the newly created user.
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

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body LoginInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 404 {string} string "Invalid credentials."
// @Failure 500 {string} string "An unexpected error occurred on the server."
// @Router /auth/sign-in [post]
// signIn is a handler for login. It returns JWT token.
func (h *Handler) signIn(c *gin.Context) {
	var input LoginInput

	if err := c.Bind(&input); err != nil {
		log.Printf("incorrect input: %v\n", err)
		c.String(http.StatusNotFound, "incorrect input")
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
