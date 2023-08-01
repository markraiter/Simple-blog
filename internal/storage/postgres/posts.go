package postgres

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/markraiter/simple-blog/models"
)

type Posts interface {
	GetAll() ([]models.Post, error)
	Filter(userID uint) ([]models.Post, error)
	Get(id uint) (*models.Post, error)
	Create(post *models.Post, userID uint) (uint, error)
	Update(id uint, input *models.UpdatePostInput) error
	Delete(id uint) error
}

type Post struct {
	db *sqlx.DB
}

func NewPost(db *sqlx.DB) *Post {
	return &Post{db: db}
}

func (s *Post) GetAll() ([]models.Post, error) {
	var posts []models.Post

	query := fmt.Sprintf("SELECT * FROM %s", postsTable)
	if err := s.db.Select(&posts, query); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *Post) Filter(userID uint) ([]models.Post, error) {
	var posts []models.Post

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", postsTable)
	if err := s.db.Select(&posts, query, userID); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *Post) Get(id uint) (*models.Post, error) {
	post := new(models.Post)

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", postsTable)

	row := s.db.QueryRow(query, id)
	if err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Body); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Post) Create(post *models.Post, userID uint) (uint, error) {
	var id uint
	query := fmt.Sprintf("INSERT INTO %s (user_id, title, body) VALUES ($1, $2, $3) RETURNING id", postsTable)

	row := s.db.QueryRow(query, userID, post.Title, post.Body)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Post) Update(id uint, input *models.UpdatePostInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET", postsTable)

	values := make(map[string]interface{})
	if input.Title != nil {
		query += " title = :title,"
		values["title"] = *input.Title
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

func (s *Post) Delete(id uint) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", postsTable)

	_, err := s.db.Exec(query, id)

	return err
}
