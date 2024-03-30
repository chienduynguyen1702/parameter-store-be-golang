package controllers

import (
	"net/http"
	"parameter-store-be/models"

	"github.com/gin-gonic/gin"
)

// Get all environments
// @Summary Get all environments
// @Description Get all environments
// @Tags Project Detail / Environment
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"environments": "environments"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to list environments"}"
// @Router /api/v1/project/{project_id}/environment [get]
func GetEnvironments(c *gin.Context) {
	projectID, exist := c.Get("project_id")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project ID from user"})
		return
	}

	var environments []models.Environment
	if err := DB.Find(&environments, "project_id = ?", projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list environments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"environments": environments})
}

// Create new environment
// @Summary Create new environment
// @Description Create new environment
// @Tags Project Detail / Environment
// @Accept json
// @Produce json
// @Param Environment body controllers.CreateEnvironment.createEnvironmentRequestBody true "Environment"
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"environment": "environment"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to create environment"}"
// @Router /api/v1/project/{project_id}/environment [post]
func CreateEnvironment(c *gin.Context) {
	projectID, exist := c.Get("project_id")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project ID from user"})
		return
	}

	type createEnvironmentRequestBody struct {
		Name string `json:"name" binding:"required"`
	}
	newEnvironmentBody := createEnvironmentRequestBody{}
	if err := c.ShouldBindJSON(&newEnvironmentBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	environment := models.Environment{
		Name:      newEnvironmentBody.Name,
		ProjectID: projectID.(uint),
	}

	if err := DB.Create(&environment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create environment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"environment": environment})
}

// Delete environment
// @Summary Delete environment
// @Description Delete environment
// @Tags Project Detail / Environment
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param environment_id path string true "Environment ID"
// @Success 200 string {string} json "{"message": "Environment deleted"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to delete environment"}"
// @Router /api/v1/project/{project_id}/environment/{environment_id} [delete]
func DeleteEnvironment(c *gin.Context) {
	projectID, exist := c.Get("project_id")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project ID from user"})
		return
	}

	environmentID := c.Param("environment_id")

	var environment models.Environment
	if err := DB.First(&environment, "id = ? AND project_id = ?", environmentID, projectID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Environment not found"})
		return
	}

	if err := DB.Delete(&environment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete environment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Environment deleted"})
}

// Update environment
// @Summary Update environment
// @Description Update environment
// @Tags Project Detail / Environment
// @Accept json
// @Produce json
// @Param Environment body controllers.UpdateEnvironment.updateEnvironmentRequestBody true "Environment"
// @Param project_id path string true "Project ID"
// @Param environment_id path string true "Environment ID"
// @Success 200 string {string} json "{"environment": "environment"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to update environment"}"
// @Router /api/v1/project/{project_id}/environment/{environment_id} [put]
func UpdateEnvironment(c *gin.Context) {
	projectID, exist := c.Get("project_id")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project ID from user"})
		return
	}

	environmentID := c.Param("environment_id")

	type updateEnvironmentRequestBody struct {
		Name string `json:"name" binding:"required"`
	}
	updateEnvironmentBody := updateEnvironmentRequestBody{}
	if err := c.ShouldBindJSON(&updateEnvironmentBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	var environment models.Environment
	if err := DB.First(&environment, "id = ? AND project_id = ?", environmentID, projectID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Environment not found"})
		return
	}

	environment.Name = updateEnvironmentBody.Name

	if err := DB.Save(&environment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update environment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"environment": environment})
}
