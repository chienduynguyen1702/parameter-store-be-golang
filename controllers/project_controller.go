package controllers

import (
	"net/http"
	"parameter-store-be/models"
	"time"

	"github.com/gin-gonic/gin"
)

// ListProjects godoc
// @Summary List projects
// @Description List projects
// @Tags Projects
// @Accept json
// @Produce json
// @Success 200 string {string} json "{"projects": "projects"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to list projects"}"
// @Router /api/v1/projects/ [get]
func ListProjects(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	userOrganizationID := user.(models.User).OrganizationID
	if userOrganizationID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization ID from user"})
		return
	}

	var projects []models.Project
	DB.Where("organization_id = ?", userOrganizationID).Find(&projects)

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

type projectBody struct {
	Name          string `gorm:"type:varchar(100);not null"`
	StartAt       time.Time
	Description   string `gorm:"type:text"`
	CurrentSprint int
	RepoURL       string `gorm:"type:varchar(100);not null"`
}

// CreateNewProject godoc
// @Summary Create new project
// @Description Create new project for organization
// @Tags Projects
// @Accept json
// @Produce json
// @Param Project body projectBody true "Project"
// @Success 200 string {string} json "{"project": "project"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to create project"}"
// @Router /api/v1/projects/ [post]
func CreateNewProject(c *gin.Context) {
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

	// Bind JSON data to projectBody struct
	var requestBody projectBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new project
	project := models.Project{
		OrganizationID: userOrganizationID,
		Name:           requestBody.Name,
		StartAt:        requestBody.StartAt,
		Description:    requestBody.Description,
		CurrentSprint:  requestBody.CurrentSprint,
		RepoURL:        requestBody.RepoURL,
	}

	// Save the new project to the database
	DB.Create(&project)

	c.JSON(http.StatusOK, gin.H{"project": project})
}

// UpdateProjectInformation godoc
// @Summary Update project information
// @Description Update project information
// @Tags Projects
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param Project body projectBody true "Project"
// @Success 200 string {string} json "{"project": "project"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to update project"}"
// @Router /api/v1/projects/{project_id} [put]
func UpdateProjectInformation(c *gin.Context) {
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

	// Bind JSON data to projectBody struct
	var requestBody projectBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	// Update project fields
	project.Name = requestBody.Name
	project.StartAt = requestBody.StartAt
	project.Description = requestBody.Description
	project.CurrentSprint = requestBody.CurrentSprint
	project.RepoURL = requestBody.RepoURL

	// Save the updated project back to the database
	DB.Save(&project)

	c.JSON(http.StatusOK, gin.H{"project": project})
}

// DeleteProject godoc
// @Summary Delete project
// @Description Delete project
// @Tags Projects
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"message": "Project deleted"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to delete project"}"
// @Router /api/v1/projects/{project_id} [delete]
func DeleteProject(c *gin.Context) {
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

	// Delete the project from the database
	DB.Delete(&project)

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted"})
}
