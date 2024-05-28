package model

import (
	"strings"

	"github.com/go-playground/validator"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"min=3,max=50" example:"username"`
	Password string `json:"password" validate:"min=8,max=50,number,upper,lower,special" example:"Password12345!"`
	Email    string `json:"email" validate:"email" example:"email@example.com"`
}

type UserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" example:"username"`
	Password string `json:"password" validate:"required,min=8,max=50,number,upper,lower,special" example:"Password12345!"`
	Email    string `json:"email" validate:"required,email" example:"email@example.com"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"email@example.com"`
	Password string `json:"password" validate:"required,min=8,max=50,number,upper,lower,special" example:"Password12345!"`
}

// ValidateContainsNumber checks if password contains at least one number
//
// Example:
//
//	ValidateContainsNumber("password12345") // true
//	ValidateContainsNumber("password") // false
func ValidateContainsNumber(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	for _, char := range password {
		if char >= '0' && char <= '9' {
			return true
		}
	}

	return false
}

// ValidateContainsUpper checks if password contains at least one uppercase letter
//
// Example:
//
//	ValidateContainsUpper("Password") // true
//	ValidateContainsUpper("password") // false
func ValidateContainsUpper(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			return true
		}
	}

	return false
}

// ValidateContainsLower checks if password contains at least one lowercase letter
//
// Example:
//
//	ValidateContainsLower("password") // true
//	ValidateContainsLower("PASSWORD") // false
func ValidateContainsLower(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	for _, char := range password {
		if char >= 'a' && char <= 'z' {
			return true
		}
	}

	return false
}

// ValidateContainsSpecial checks if password contains at least one special character
//
// Example:
//
//	ValidateContainsSpecial("password!") // true
//	ValidateContainsSpecial("password") // false
func ValidateContainsSpecial(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	specialChars := "!@#$%^&*()-_=+[]{}|;:'\",.<>/?"
	for _, char := range password {
		if strings.ContainsRune(specialChars, char) {
			return true
		}
	}

	return false
}
