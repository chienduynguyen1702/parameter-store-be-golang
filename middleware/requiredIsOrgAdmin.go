package middleware

import (
	"net/http"
	"parameter-store-be/models"

	"github.com/gin-gonic/gin"
)

func RequiredIsOrgAdmin(c *gin.Context) {
	// get user from context
	userInContext, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	user := userInContext.(models.User)
	// check if user is organization admin
	if user.IsOrganizationAdmin {
		c.Next()
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not an organization admin"})
}
