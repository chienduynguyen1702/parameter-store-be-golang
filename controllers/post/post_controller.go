package controllers

import (
	"net/http"
	main "vcs_backend/gorm/controllers"
	"vcs_backend/gorm/models"

	"github.com/gin-gonic/gin"
)

func PostController(c *gin.Context) {
	var posts []models.Post
	main.DB.Find(&posts)
	// fmt.Println(posts)
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}
