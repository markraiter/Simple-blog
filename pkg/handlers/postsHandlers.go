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
func GetPostByID(db *gorm.DB) echo.HandlerFunc {
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

// Creating post
func CreatePost(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		post := new(models.Post)

		if err := c.Bind(post); err != nil {
			return c.String(http.StatusBadRequest, "invalid post data")
		}

		if err := db.Create(post).Error; err != nil {
			return c.String(http.StatusInternalServerError, "error creating post")
		}

		return c.JSON(http.StatusOK, post)
	}
}

// Updating post
func UpdatePost(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		postID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid post id")
		}

		post := new(models.Post)
		
		if err := db.First(post, postID).Error; err != nil {
			return c.String(http.StatusNotFound, "post not found")
		}

		if err := c.Bind(post); err != nil {
			return c.String(http.StatusBadRequest, "invalid post data")
		}

		if err := db.Save(post).Error; err != nil {
			return c.String(http.StatusInternalServerError, "error updating post")
		}

		return c.JSON(http.StatusOK, post)
	}
}

// Deleting post
func DeletePost(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		postID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid posy id")
		}

		post := new(models.Post)

		if err := db.First(post, postID).Error; err != nil {
			return c.String(http.StatusNotFound, "post not found")
		}

		if err := db.Delete(post).Error; err != nil {
			return c.String(http.StatusInternalServerError, "error deleting post")
		}

		return c.NoContent(http.StatusNoContent)
	}
}