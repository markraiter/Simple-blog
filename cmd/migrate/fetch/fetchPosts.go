package fetch

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/markraiter/simple-blog/internal/initializers"
	"github.com/markraiter/simple-blog/internal/models"
)

func FetchPosts() {
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

	for i := range p {
		if err := initializers.DB.Create(&p[i]); err != nil {
			log.Printf("error writing posts to database: %v", err)
		}
	}
}