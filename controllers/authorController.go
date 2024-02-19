package controllers

import (
	"net/http"
	"vcs_backend/gorm/models"

	"github.com/gin-gonic/gin"
)

func AuthorController(c *gin.Context) {
	var authors []models.Author
	db.Find(&authors)
	// fmt.Println(authors)
	c.JSON(http.StatusOK, gin.H{
		"authors": authors,
	})
}
