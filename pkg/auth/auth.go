package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/markraiter/simple-blog/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	_ "github.com/markraiter/simple-blog/docs"
)

func generateToken(u *models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":       u.ID,
		"email":    u.Email,
		"password": u.Password,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// @Summary Create a new user account
// @Tags Authentification
// @Description Register a new user account
// @Accept json
// @Produce json
// @Param user body models.User true "User object"
// @Success 200 {object} models.User
// @Failure 400 {object} error "Bad request"
// @Failure 500 {object} error "Internal server error"
// @Router /registration [post]
func Register(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := new(models.User)

		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		user.Password = string(hashedPassword)

		if err := db.Create(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, user)
	}
}

// @Summary Login
// @Tags Authentification
// @Description Logs in user with email and password
// @Accept  json
// @Produce  json
// @Param user body models.User true "User object"
// @Success 200 {string} string "token"
// @Failure 400 {object} error "Bad request"
// @Failure 401 {object} error "Unauthorized"
// @Failure 500 {object} error "Internal server error"
// @Router /login [post]
func Login(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := new(models.User)

		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		var dbUser models.User

		if err := db.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid input",
			})
		}

		token, err := generateToken(&dbUser)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": token,
		})
	}
}
