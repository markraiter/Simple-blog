package postgres

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/markraiter/simple-blog/models"
)

type Comment struct {
	db *sqlx.DB
}

func NewComment(db *sqlx.DB) *Comment {
	return &Comment{db: db}
}

func (s *Comment) GetAll() ([]models.Comment, error) {
	var comments []models.Comment

	query := fmt.Sprintf("SELECT * FROM %s", commentsTable)

	if err := s.db.Select(&comments, query); err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *Comment) FilterByPost(postID uint) ([]models.Comment, error) {
	var comments []models.Comment

	query := fmt.Sprintf("SELECT * FROM %s WHERE post_id=$1", commentsTable)
	if err := s.db.Select(&comments, query, postID); err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *Comment) FilterByUser(userID uint) ([]models.Comment, error) {
	var comments []models.Comment

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", commentsTable)
	if err := s.db.Select(&comments, query, userID); err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *Comment) Get(id uint) (*models.Comment, error) {
	comment := new(models.Comment)

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", commentsTable)

	row := s.db.QueryRow(query, id)
	if err := row.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Email, &comment.Body); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *Comment) Create(comment *models.Comment, userID, postID uint) (uint, error) {
	var id uint
	query := fmt.Sprintf("INSERT INTO %s (user_id, post_id, email, body) VALUES ($1, $2, $3, $4) RETURNING id", commentsTable)

	row := s.db.QueryRow(query, userID, postID, comment.Email, comment.Body)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Comment) Update(id uint, input *models.UpdateCommentInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET", commentsTable)

	values := make(map[string]interface{})
	if input.Email != nil {
		query += " email = :email,"
		values["email"] = *input.Email
	}
	if input.Body != nil {
		query += " body = :body,"
		values["body"] = *input.Body
	}

	query = strings.TrimSuffix(query, ",")

	query += " WHERE id = :id"

	values["id"] = id

	_, err := s.db.NamedExec(query, values)
	if err != nil {
		return err
	}

	return nil
}

func (s *Comment) Delete(id uint) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", commentsTable)

	_, err := s.db.Exec(query, id)

	return err
}
