package controllers

import (
	"net/http"

	main "vcs_backend/gorm/controllers"
	"vcs_backend/gorm/models"

	"github.com/gin-gonic/gin"
)

// CREATE
// func CreateNewPost(c *gin.Context) {
// 	var post models.Post
// 	post.Title = c.Query("title")
// 	post.Body = c.Request.Body.String().TrimSpace()
// 	if err := main.DB.Create(&post); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err})
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "success",
// 			"data":    post,
// 		})
// 	}
// }

// READ
func GetPosts(c *gin.Context) {
	authorId := c.Query("author-id")
	if authorId != "" {
		GetPostsByAuthorID(c)
	} else {
		GetAllPosts(c)
	}
}

func GetAllPosts(c *gin.Context) {
	var posts []models.Post
	if err := main.DB.Find(&posts).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func GetPostsByAuthorID(c *gin.Context) {
	authorID := c.Query("author-id")

	var posts []models.Post
	if err := main.DB.
		Joins("JOIN author_posts ON posts.id = author_posts.post_id").
		Where("author_posts.author_id = ?", authorID).
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}
