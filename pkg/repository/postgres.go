package repository

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

const (
	usersTable    = "users"
	postsTable    = "posts"
	commentsTable = "comments"
)

type Config struct {
	Driver   string
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open(cfg.Driver, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		log.Fatalf("error connecting to database: %s\n", err.Error())
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
