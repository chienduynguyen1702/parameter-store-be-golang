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
// @Security ApiKeyAuth
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
	from := c.Query("from")
	to := c.Query("to")
	// Retrieve tracking from the database using the project ID
	var agentLogs []models.AgentLog
	var projectLogs []models.ProjectLog
	if from != "" && to != "" {

		from := startOfDay(c.Query("from"))
		to := endOfDay(c.Query("to"))

		DB.Preload("Agent").Where("project_id = ? AND created_at BETWEEN ? AND ?", projectID, from, to).Find(&agentLogs)
		DB.Preload("User").Where("project_id = ? AND created_at BETWEEN ? AND ?", projectID, from, to).Find(&projectLogs)
	} else {
		DB.Preload("Agent").Where("project_id = ?", projectID).Find(&agentLogs)
		DB.Preload("User").Where("project_id = ?", projectID).Find(&projectLogs)
	}
	// fmt.Println(agentLogs)
	// fmt.Println(projectLogs)
	// Combine agentLogs and projectLogs
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"agent_logs":   agentLogs,
			"project_logs": projectLogs,
		},
	})
}

func startOfDay(from string) string {
	return from + " 00:00:00"
}
func endOfDay(to string) string {
	return to + " 23:59:59"
}
