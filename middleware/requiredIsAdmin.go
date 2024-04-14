package middleware

import (
	"net/http"
	"parameter-store-be/controllers"
	"parameter-store-be/models"

	"github.com/gin-gonic/gin"
)

func RequiredIsAdmin(c *gin.Context) {
	// get user from context
	userInContext, exists := c.Get("user")
	if !exists {
		// log.Println("Failed to get user from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	user := userInContext.(models.User)
	// check if user is organization admin
	if user.IsOrganizationAdmin {
		// log.Println("User is organization admin")
		c.Next()
		return
	}

	// get project_id from path
	project_id := c.Param("project_id")
	if project_id == "0" {
		// log.Println("Failed to get project ID from user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project ID from user"})
		return
	}
	// check if user belongs to the project
	var upr models.UserRoleProject
	if err := controllers.DB.Where("user_id = ? AND project_id = ?", user.ID, project_id).First(&upr).Error; err != nil {
		// log.Println("Failed to get user role project")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User does not belong to the project"})
		return
	}
	if upr.Role.Name != "Project Admin" {
		// log.Println("User is not an admin, please contact the project admin to perform this action")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not an admin, please contact the project admin to perform this action"})
		return
	}
	c.Next()
}
