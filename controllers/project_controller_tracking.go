package controllers

import (
	"net/http"

	"parameter-store-be/models"

	"github.com/gin-gonic/gin"
)

// GetProjectTracking is a function to get project tracking
// @Summary Get project tracking
// @Description Get project tracking
// @Tags Project Detail / Tracking
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"tracking": "tracking"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get project tracking"}"
// @Router /api/v1/projects/{project_id}/tracking [get]
func GetProjectTracking(c *gin.Context) {
	// Retrieve user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	// Type assertion to extract organization ID
	userOrganizationID := user.(models.User).OrganizationID
	if userOrganizationID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization ID from user"})
		return
	}

	// Retrieve project ID from the URL
	projectID := c.Param("project_id")

	// Retrieve project from the database using the project ID
	var project models.Project
	result := DB.First(&project, projectID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project"})
		return
	}

	// Retrieve tracking from the database using the project ID
	var logs []models.AgentLog
	DB.Preload("Agent").Where("project_id = ?", projectID).Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"tracking": logs,
		},
	})
}
