package middleware

import (
	"net/http"
	"parameter-store-be/models"

	"github.com/gin-gonic/gin"
)

func RequiredBelong(c *gin.Context) {
	// get user from context
	userInContext, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	// get project_id from path
	project_id := c.Param("project_id")
	if project_id == "0" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project ID from user"})
		return
	}
	user := userInContext.(models.User)
	// check if user belongs to the project in model user-project-role
	if user.ID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

}
