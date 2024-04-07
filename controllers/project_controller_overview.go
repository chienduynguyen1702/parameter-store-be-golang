package controllers

import (
	"net/http"
	"parameter-store-be/models"
	"time"

	"github.com/gin-gonic/gin"
)

// GetProjectDetail godoc
// @Summary Get project overview
// @Description Get project overview
// @Tags Project Detail / Overview
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"project": "project"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get project detail"}"
// @Router /api/v1/projects/{project_id}/overview [get]
func GetProjectOverView(c *gin.Context) {
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
	// Retrieve users and their roles in the given project
	var upr []models.UserProjectRole
	DB.Preload("User").Preload("Role").Where("project_id = ?", projectID).Find(&upr)

	// // user current in project, by UserProjectRole table
	// var stagesInProject []models.Stage
	// DB.Model(&project).Association("Stages").Find(&stagesInProject)
	// project.Stages = stagesInProject
	// // user current in project, by UserProjectRole table
	// var environmentsInProject []models.Environment
	// DB.Model(&project).Association("Environments").Find(&environmentsInProject)
	// project.Environments = environmentsInProject
	// // user current in project, by UserProjectRole table
	// var agentsInProject []models.Agent
	// DB.Model(&project).Association("Agents").Find(&agentsInProject)
	// project.Agents = agentsInProject
	// help me this
	type UserRoleInProject struct {
		UserID   uint   `json:"id"`
		UserName string `json:"name"`
		RoleName string `json:"role"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		// LastLogIn time.Time `json:"last_login"`
	}
	var userRoleInProject []UserRoleInProject
	for _, upr := range upr {
		userRoleInProject = append(userRoleInProject, UserRoleInProject{
			UserID:   upr.User.ID,
			UserName: upr.User.Username,
			RoleName: upr.Role.Name,
			Email:    upr.User.Email,
			Phone:    upr.User.Phone,
			// LastLogIn: upr.User.LastLogIn,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"overview": project,
		"users":    userRoleInProject,
	})
}

type projectBody struct {
	Name          string `gorm:"type:varchar(100);not null"`
	StartAt       time.Time
	Description   string `gorm:"type:text"`
	CurrentSprint string
	RepoURL       string `gorm:"type:varchar(100);not null"`
}

// UpdateProjectInformation godoc
// @Summary Update project information
// @Description Update project information
// @Tags Project Detail / Overview
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param Project body projectBody true "Project"
// @Success 200 string {string} json "{"project": "project"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to update project"}"
// @Router /api/v1/projects/{project_id}/overview [put]
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
