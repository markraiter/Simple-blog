package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/markraiter/simple-blog/models"
)

type PostPostgres struct {
	db *sqlx.DB
}

func NewPostPostgres(db *sqlx.DB) *PostPostgres {
	return &PostPostgres{
		db: db,
	}
}

func (r *PostPostgres) Create(userID int, post models.Post) (int, error) {
	var id int
	createPostQuery := fmt.Sprintf("INSERT INTO %s (userId, title, body) VALUES ($1, $2, $3) RETURNING id", postsTable)

	row := r.db.QueryRow(createPostQuery, post.Title, post.Body)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
