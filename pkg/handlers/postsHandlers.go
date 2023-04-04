package handlers

import (
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/markraiter/simple-blog/internal/models"
	"gorm.io/gorm"
	_ "github.com/markraiter/simple-blog/docs"
)

// @Summary Get all posts
// @Description Get all posts from the database.
// @Tags Posts
// @Accept  json
// @Produce  json
// @Produce  xml
// @Security ApiKeyAuth
// @Success 200 {array} models.Post
// @Failure 500 {string} error fetching post
// @Router /api/v1/posts [get]
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

// @Summary Get post by ID
// @Description Get particular post from the database by unique ID
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Security ApiKeyAuth
// @Success 200 {object} models.Post
// @Failure 400 {string} invalid post ID
// @Failure 404 {string} post not found
// @Router /api/v1/posts/{id} [get] 
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

// @Summary Create post
// @Description Create new post in the database
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body models.Post true "Post data"
// @Security ApiKeyAuth
// @Success 201 {object} models.Post
// @Failure 400 {string} invalid post data
// @Failure 500 {string} error creating post
// @Router /api/v1/posts [post]
func CreatePost(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		post := new(models.Post)

		if err := c.Bind(post); err != nil {
			return c.String(http.StatusBadRequest, "invalid post data")
		}

		if err := db.Create(post).Error; err != nil {
			return c.String(http.StatusInternalServerError, "error creating post")
		}

		return c.JSON(http.StatusCreated, post)
	}
}

// @Summary Update post
// @Description Update information in the particular post from the database by unique ID
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body models.Post true "Post data"
// @Security ApiKeyAuth
// @Success 200 {object} models.Post
// @Failure 400 {string} invalid post data
// @Failure 404 {string} post not found
// @Failure 500 {string} error updating post
// @Router /api/v1/posts/{id} [put]
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

// @Summary Delete post
// @Description Delete particular post from the database by unique ID
// @Tags Posts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 204 {string} post deleted successfully
// @Failure 400 {string} invalid post ID
// @Failure 404 {string} post not found
// @Failure 500 {string} error deleting post
// @Router /api/v1/posts/{id} [delete]
func DeletePost(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		postID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid post id")
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
