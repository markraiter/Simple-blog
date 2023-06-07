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
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createPostQuery := fmt.Sprintf("INSERT INTO %s (title, body) VALUES ($1, $2)", postsTable)

	row := tx.QueryRow(createPostQuery, post.Title, post.Body)
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}
