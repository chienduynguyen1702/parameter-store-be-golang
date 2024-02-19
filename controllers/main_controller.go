package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

// SetDB sets the db object
func SetDB(database *gorm.DB) {
	db = database
}

func MainController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "this is the main page",
	})
}
