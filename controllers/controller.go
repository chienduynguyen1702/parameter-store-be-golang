package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

// SetDB sets the db object
func SetDB(database *gorm.DB) {
	DB = database
}

func MainController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "this is first index api ",
	})
}
