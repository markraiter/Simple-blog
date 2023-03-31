package handlers

import (
	"encoding/xml"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markraiter/simple-blog/internal/models"
	"gorm.io/gorm"
)

// Getting all comments
func GetComments(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		responseType := c.QueryParam("type")
		if responseType != "xml" {
			responseType = "json"
		}
		
		comments := []models.Comment{}

		if err := db.Find(&comments).Error; err != nil {
			return c.String(http.StatusInternalServerError, "error retreiving comments")
		}

		if responseType == "xml" {
			xmlResponse, err := xml.MarshalIndent(comments, "", " ")
			if err != nil {
				return c.String(http.StatusInternalServerError, "error encoding xml")
			}

			c.Response().Header().Set("content-type", "application/xml")

			return c.Blob(http.StatusOK, "application/xml", xmlResponse)
			
		} else {
			c.Response().Header().Set("content-type", "application/json")

			return c.JSON(http.StatusOK, comments)
		}
	}
}