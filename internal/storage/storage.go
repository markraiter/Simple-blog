package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/markraiter/simple-blog/internal/storage/postgres"
)

type Storage struct {
	postgres.Authentication
	postgres.Posts
	postgres.Comments
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		postgres.NewAuth(db),
		postgres.NewPost(db),
		postgres.NewComment(db),
	}
}
