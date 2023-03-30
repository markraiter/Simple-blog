package fetch

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/markraiter/simple-blog/internal/initializers"
	"github.com/markraiter/simple-blog/internal/models"
)

func FetchComments() {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/comments?postId=7")
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

	for i := range c {
		if err := initializers.DB.Create(&c[i]); err != nil {
			log.Printf("error writing comments to database: %v", err)
		}
	}
}