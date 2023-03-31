package migrate

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/markraiter/simple-blog/internal/models"
	"gorm.io/gorm"
)

func fetchPosts() ([]models.Post, error) {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts?userId=7")
	if err != nil {
		log.Printf("errror fetching posts: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
	}

	var p []models.Post
	parseError := json.Unmarshal(body, &p)
	if parseError != nil {
		log.Println(parseError.Error())
	}

	return p, nil
}

func fetchComments(postId int) ([]models.Comment, error) {
	resp, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/comments?postId=%d", postId))
	if err != nil {
		log.Printf("error fetching comments: %s", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
	}

	var c []models.Comment

	parseError := json.Unmarshal(body, &c)
	if parseError != nil {
		log.Println(parseError.Error())
	}

	return c, nil
}

func WriteToDB(db * gorm.DB) error {
	var comments []models.Comment
	posts, err := fetchPosts()
	if err != nil {
		log.Printf("error fetching posts: %s", err)
	}

	for _, post := range posts {
		c, err := fetchComments(int(post.ID))
		if err != nil {
			log.Printf("error fetching comments: %s", err)
		}
		comments = append(comments, c...)
	}

	for i := range posts {
		if err := db.Create(&posts[i]); err != nil {
			log.Printf("error writing posts to database: %v", err)
		}
	}

	for v := range comments {
		if err := db.Create(&comments[v]); err != nil {
			log.Printf("error writing comments to database: %v", err)
		}
	}

	return nil
}