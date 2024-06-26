package controllers

import (
	"net/http"
	"parameter-store-be/models"

	"github.com/gin-gonic/gin"
)

// GetProjectAllInfo is a function to get all project info
// @Summary Get all project info
// @Description Get all project info
// @Tags Project
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"projects": "projects"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Security ApiKeyAuth
// @Failure 500 string {string} json "{"error": "Failed to get project info"}"
// @Router /api/v1/projects/{project_id} [get]
func GetProjectAllInfo(c *gin.Context) {
	project_id := c.Param("project_id")
	if project_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	var project models.Project
	// preload stage, environment, agent, parameter in project
	// result := DB.First(&project, project_id)
	// if result.Error != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project"})
	// 	return
	// }

	result := DB.Preload("Stages").
		Preload("Versions").
		Preload("Versions.Parameters").
		Preload("Versions.Parameters.Stage").
		Preload("Versions.Parameters.Environment").
		Preload("LatestVersion").
		Preload("LatestVersion.Parameters").
		Preload("LatestVersion.Parameters.Stage").
		Preload("LatestVersion.Parameters.Environment").
		Preload("Environments").
		Preload("Agents").
		Preload("Parameters").
		Preload("UserRoles").
		Preload("UserRoles.User"). // Preload User association
		Preload("UserRoles.Role"). // Preload Role association
		Preload("Logs").
		Preload("Workflows").
		Preload("Workflows.Logs").
		First(&project, project_id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"projects": project})
}
