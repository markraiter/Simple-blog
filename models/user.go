package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint   `json:"-" db:"id"`
	Email        string `json:"email" binding:"required"`
	Password     string `json:"password" db:"-" binding:"required"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
}

// generatePasswordHash function generates hash from input password
func generatePasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing the password: %w", err)
	}

	return string(hashedPassword), nil
}

// ComparePassword function compares input password with hashed password in the database
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}

// BeforeCreate function generates hash from password if len(password) > 0
func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		hash, err := generatePasswordHash(u.Password)
		if err != nil {
			return err
		}

		u.PasswordHash = hash
	}

	return nil
}
