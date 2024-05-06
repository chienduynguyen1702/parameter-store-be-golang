package controllers

import (
	"log"
	"net/http"
	"parameter-store-be/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Get all environments in a project
// @Summary Get all environments in a project
// @Description Get all environments in a project
// @Tags Project Detail / Environments
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"environments": "environments"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to list environments"}"
// @Router /api/v1/projects/{project_id}/environments [get]
func GetListEnvironmentInProject(c *gin.Context) {
	// Retrieve project ID from the URL
	projectID := c.Param("project_id")

	// Retrieve environments from the database using the project ID
	var environments []models.Environment
	result := DB.Where("project_id = ? AND is_archived = ? ", projectID, false).Find(&environments)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve environments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"environments": environments,
		},
	})
}

// Get a environment in a project
// @Summary Get a environment in a project
// @Description Get a environment in a project
// @Tags Project Detail / Environments
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param environment_id path string true "Environment ID"
// @Success 200 string {string} json "{"environment": "environment"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get environment"}"
// @Router /api/v1/projects/{project_id}/environments/{environment_id} [get]
func GetEnvironmentInProject(c *gin.Context) {
	// Retrieve project ID from the URL
	projectID := c.Param("project_id")
	// Retrieve environment ID from the URL
	environmentID := c.Param("environment_id")

	// Retrieve environment from the database using the project ID and environment ID
	var environment models.Environment
	result := DB.Where("project_id = ? AND id = ? AND is_archived = ? ", projectID, environmentID, false).First(&environment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve environment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"environment": environment,
		},
	})
}

type environmentRequestBody struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// Create a environment in a project
// @Summary Create a environment in a project
// @Description Create a environment in a project
// @Tags Project Detail / Environments
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param request body controllers.environmentRequestBody true "Environment creation request"
// @Success 201 string {string} json "{"message": "Environment created successfully", "environment": "environment"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to create environment"}"
// @Router /api/v1/projects/{project_id}/environments [post]
func CreateEnvironmentInProject(c *gin.Context) {
	// Retrieve project ID from the URL
	projectID := c.Param("project_id")
	projectIDInt, err := strconv.Atoi(projectID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	// Parse project ID to uint
	projectIDUint := uint(projectIDInt)

	// Retrieve the request body
	r := environmentRequestBody{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create a new environment in the database
	newEnvironment := models.Environment{
		Name:        r.Name,
		Description: r.Description,
		Color:       r.Color,
		ProjectID:   projectIDUint,
	}
	if err := DB.Create(&newEnvironment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create environment"})
		return
	}

	// get username from context
	user, exist := c.Get("user")
	if !exist {
		log.Println("Failed to get user from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	// Type assertion to extract username
	uID := user.(models.User).ID

	projectLogByUser(newEnvironment.ProjectID, "Updated Environment", "Environment is updated", 200, time.Since(time.Now()), uID)
	c.JSON(http.StatusCreated, gin.H{
		"message":     "Environment created successfully",
		"environment": newEnvironment,
	})
}

// Update a environment in a project
// @Summary Update a environment in a project
// @Description Update a environment in a project
// @Tags Project Detail / Environments
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param environment_id path string true "Environment ID"
// @Param request body controllers.environmentRequestBody true "Environment update request"
// @Success 200 string {string} json "{"message": "Environment updated successfully", "environment": "environment"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to update environment"}"
// @Router /api/v1/projects/{project_id}/environments/{environment_id} [put]
func UpdateEnvironmentInProject(c *gin.Context) {
	// Retrieve project ID from the URL
	projectID := c.Param("project_id")
	// Retrieve environment ID from the URL
	environmentID := c.Param("environment_id")

	// Retrieve the request body
	r := environmentRequestBody{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the environment in the database using the project ID and environment ID
	var environment models.Environment
	result := DB.Where("project_id = ? AND id = ?", projectID, environmentID).First(&environment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve environment"})
		return
	}
	environment.Name = r.Name
	environment.Description = r.Description
	environment.Color = r.Color
	if err := DB.Save(&environment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update environment"})
		return
	}

	// get username from context
	user, exist := c.Get("user")
	if !exist {
		log.Println("Failed to get user from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	// Type assertion to extract username
	uID := user.(models.User).ID

	projectLogByUser(environment.ProjectID, "Updated Environment", "Environment is updated", 200, time.Since(time.Now()), uID)
	c.JSON(http.StatusOK, gin.H{
		"message":     "Environment updated successfully",
		"environment": environment,
	})
}

// Get all archived environments in a project
// @Summary Get all archived environments in a project
// @Description Get all archived environments in a project
// @Tags Project Detail / Environments
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"environments": "environments"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to list archived environments"}"
// @Router /api/v1/projects/{project_id}/environments/archived [get]
func GetListArchivedEnvironmentInProject(c *gin.Context) {
	// Retrieve project ID from the URL
	projectID := c.Param("project_id")

	// Retrieve archived environments from the database using the project ID
	var environments []models.Environment
	result := DB.Where("project_id = ? AND is_archived = ?", projectID, true).Find(&environments)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve archived environments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"environments": environments,
		},
	})
}

// Archive a environment in a project
// @Summary Archive a environment in a project
// @Description Archive a environment in a project
// @Tags Project Detail / Environments
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param environment_id path string true "Environment ID"
// @Success 200 string {string} json "{"message": "Environment archived successfully"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to archive environment"}"
// @Router /api/v1/projects/{project_id}/environments/{environment_id}/archive [patch]
func ArchiveEnvironmentInProject(c *gin.Context) {
	// Retrieve project ID from the URL
	projectID := c.Param("project_id")
	uint64ProjectID, err := strconv.ParseUint(projectID, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Retrieve environment ID from the URL
	environmentID := c.Param("environment_id")

	// Archive the environment in the database using the project ID and environment ID
	var environment models.Environment
	result := DB.Where("project_id = ? AND id = ?", projectID, environmentID).First(&environment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve environment"})
		return
	}
	if environment.IsArchived {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Environment is already archived"})
		return
	}

	// get username from context
	user, exist := c.Get("user")
	if !exist {
		log.Println("Failed to get user from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	// Type assertion to extract username
	username := user.(models.User).Username
	uID := user.(models.User).ID

	environment.IsArchived = true
	environment.ArchivedAt = time.Now()
	environment.ArchivedBy = username
	if err := DB.Save(&environment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to archive environment"})
		return
	}
	// log to project log
	projectLogByUser(uint(uint64ProjectID), "Archived Environment", "Environment is archived", 200, time.Since(time.Now()), uID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Environment archived successfully",
	})
}

// Unarchive a environment in a project
// @Summary Unarchive a environment in a project
// @Description Unarchive a environment in a project
// @Tags Project Detail / Environments
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param environment_id path string true "Environment ID"
// @Success 200 string {string} json "{"message": "Environment unarchived successfully"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to unarchive environment"}"
// @Router /api/v1/projects/{project_id}/environments/{environment_id}/unarchive [patch]
func UnarchiveEnvironmentInProject(c *gin.Context) {
	// Retrieve project ID from the URL
	projectID := c.Param("project_id")
	// Retrieve environment ID from the URL
	environmentID := c.Param("environment_id")

	// Unarchive the environment in the database using the project ID and environment ID
	var environment models.Environment
	result := DB.Where("project_id = ? AND id = ?", projectID, environmentID).First(&environment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve environment"})
		return
	}
	if !environment.IsArchived {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Environment is already unarchived"})
		return
	}

	environment.IsArchived = false
	if err := DB.Save(&environment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unarchive environment"})
		return
	}

	// get username from context
	user, exist := c.Get("user")
	if !exist {
		log.Println("Failed to get user from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	// Type assertion to extract username
	uID := user.(models.User).ID
	// log to project log
	projectLogByUser(environment.ProjectID, "Unarchived Environment", "Environment is unarchived", 200, time.Since(time.Now()), uID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Environment unarchived successfully",
	})
}
