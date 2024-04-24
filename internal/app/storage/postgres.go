package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/markraiter/simple-blog/config"
)

type Storage struct {
	PostgresDB *sql.DB
}

func New(cfg config.Postgres) *Storage {
	entryString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Database, cfg.Password, cfg.SSLMode)
	db, err := sql.Open(cfg.Driver, entryString)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		return nil
	}

	return &Storage{PostgresDB: db}
}
