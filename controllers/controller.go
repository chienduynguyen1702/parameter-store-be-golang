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

// MainController godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func MainController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "this is first index api ",
	})
}
