package postgres

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/markraiter/simple-blog/models"
)

type Authentication interface {
	Create(user *models.User) (uint, error)
	GetEmail(email string) string
	GetUserByEmail(email, password string) (*models.User, error)
	GenerateToken(email, password string) (string, error)
}

type Auth struct {
	db *sqlx.DB
}

func NewAuth(db *sqlx.DB) *Auth {
	return &Auth{db: db}
}

func (s *Auth) Create(user *models.User) (uint, error) {
	var id uint

	if err := user.BeforeCreate(); err != nil {
		return 0, err
	}

	query := fmt.Sprintf("INSERT INTO %s (email, password_hash) VALUES ($1, $2) RETURNING id", usersTable)
	row := s.db.QueryRow(query, user.Email, user.PasswordHash)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Auth) GetEmail(email string) string {
	var userEmail string

	query := fmt.Sprintf("SELECT email FROM %s WHERE email=$1", usersTable)

	if err := s.db.Get(&userEmail, query, email); err != nil {
		return ""
	}

	return userEmail
}

func (s *Auth) GetUserByEmail(email, password string) (*models.User, error) {
	user := new(models.User)

	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1", usersTable)

	row := s.db.QueryRow(query, email)
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		return nil, err
	}

	if !user.ComparePassword(password) {
		return nil, errors.New("password doesn't match")
	}

	return user, nil
}

func (s *Auth) GenerateToken(email, password string) (string, error) {
	user, err := s.GetUserByEmail(email, password)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"username": user.Email,
		"exp":      time.Now().Add(12 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
