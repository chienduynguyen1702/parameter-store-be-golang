package controllers

import (
	"net/http"
	"parameter-store-be/models"
	"strconv"
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
	var urp []models.UserRoleProject
	DB.Preload("User").Preload("Role").Where("project_id = ?", projectID).Find(&urp)

	// // user current in project, by UserRoleProject table
	// var stagesInProject []models.Stage
	// DB.Model(&project).Association("Stages").Find(&stagesInProject)
	// project.Stages = stagesInProject
	// // user current in project, by UserRoleProject table
	// var environmentsInProject []models.Environment
	// DB.Model(&project).Association("Environments").Find(&environmentsInProject)
	// project.Environments = environmentsInProject
	// // user current in project, by UserRoleProject table
	// var agentsInProject []models.Agent
	// DB.Model(&project).Association("Agents").Find(&agentsInProject)
	// project.Agents = agentsInProject
	// help me this
	type UserRoleInProject struct {
		UserID   uint   `json:"id"`
		UserName string `json:"username"`
		RoleName string `json:"role"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		// LastLogIn time.Time `json:"last_login"`
	}
	var userRoleInProject []UserRoleInProject
	for _, urp := range urp {
		userRoleInProject = append(userRoleInProject, UserRoleInProject{
			UserID:   urp.User.ID,
			UserName: urp.User.Username,
			RoleName: urp.Role.Name,
			Email:    urp.User.Email,
			Phone:    urp.User.Phone,
			// LastLogIn: urp.User.LastLogIn,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"overview": project,
		"users":    userRoleInProject,
	})
}

type projectBody struct {
	Name          string    `gorm:"type:varchar(100);not null" json:"name"`
	StartAt       time.Time `json:"start_at"`
	Description   string    `gorm:"type:text" json:"description"`
	CurrentSprint string    `gorm:"type:varchar(100)" json:"current_sprint"`
	RepoURL       string    `gorm:"type:varchar(100);not null" json:"repo_url"`
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

// AddUserToProject godoc
// @Summary Add user to project include role
// @Description Add user to project include role
// @Tags Project Detail / Overview
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param UserRoleProject body controllers.AddUserToProject.UserRoleProjectBody true "UserRoleProject"
// @Success 200 string {string} json "{"message": "User added to project"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to add user to project"}"
// @Router /api/v1/projects/{project_id}/overview/add-user [post]
func AddUserToProject(c *gin.Context) {
	// Bind JSON data to UserRoleProject struct
	type UserRoleProjectBody struct {
		UserID uint ` json:"user_id"`
		RoleID uint ` json:"role_id"`
	}
	var urpb UserRoleProjectBody
	if err := c.ShouldBindJSON(&urpb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve project ID from the URL
	projectID := c.Param("project_id")
	parsedProjectID, err := strconv.ParseUint(projectID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var urp models.UserRoleProject
	// Check if the user is already in the project
	result := DB.Where("user_id = ? AND project_id = ?", urpb.UserID, parsedProjectID).First(&urp)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already in the project"})
		return
	}

	// Create a new user to project relationship
	urp = models.UserRoleProject{
		UserID:    urpb.UserID,
		ProjectID: uint(parsedProjectID),
		RoleID:    urpb.RoleID,
	}
	// Save the new user to project relationship to the database
	DB.Create(&urp)

	c.JSON(http.StatusOK, gin.H{"message": "User added to project"})
}
