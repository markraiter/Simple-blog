package models

type User struct {
	ID       string `json:"-" db:"id"`
	Email    string `json:"email" xml:"email"`
	Password string `json:"password" xml:"password"`
}
