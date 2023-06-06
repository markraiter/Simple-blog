package repository

import (
	"database/sql"
	"fmt"

	"github.com/markraiter/simple-blog/models"
)

type PostMySQL struct {
	db *sql.DB
}

func NewPostMySQL(db *sql.DB) *PostMySQL {
	return &PostMySQL{
		db: db,
	}
}

func (r *PostMySQL) Create(userID int, post models.Post) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createPostQuery := fmt.Sprintf("INSERT INTO %s (title, body) VALUES (?, ?)", postsTable)
	row := tx.QueryRow(createPostQuery, post.Title, post.Body)

	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}
