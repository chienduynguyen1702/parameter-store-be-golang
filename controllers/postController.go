package controllers

import (
	"fmt"
	"net/http"
	"vcs_backend/gorm/models"

	"github.com/gin-gonic/gin"
)

func PostController(c *gin.Context) {
	var posts []models.Post
	db.Find(&posts)
	fmt.Println(posts)
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}
