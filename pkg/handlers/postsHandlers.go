package handlers

import (
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/markraiter/simple-blog/internal/models"
	"gorm.io/gorm"
)

// Getting all posts
func GetPosts(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		responseType := c.QueryParam("type")
		if responseType != "xml" {
			responseType = "json"
		}

		posts := []models.Post{}

		if err := db.Find(&posts).Error; err != nil {
			return c.String(http.StatusInternalServerError, "error retreiving posts")
		}

		if responseType == "xml" {
			xmlResponse, err := xml.MarshalIndent(posts, "", " ")
			if err != nil {
				return c.String(http.StatusInternalServerError, "error encoding xml")
			}

			c.Response().Header().Set("content-type", "application/xml")
			return c.Blob(http.StatusOK, "application/xml", xmlResponse)
		} else {
			c.Response().Header().Set("content-type", "application/json")
			return c.JSON(http.StatusOK, posts)
		}
	}
}

// Getting post by ID 
func GetPost(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		postID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid post ID")
		}

		post := new(models.Post)
		if err := db.First(post, postID).Error; err != nil {
			return c.String(http.StatusNotFound, "post not found")
		}

		return c.JSON(http.StatusOK, post)
	}
}