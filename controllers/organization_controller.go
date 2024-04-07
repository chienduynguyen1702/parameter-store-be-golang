package controllers

import (
	"fmt"
	"net/http"
	"parameter-store-be/models"
	"time"

	"github.com/gin-gonic/gin"
)

// GetOrganizationInformation godoc
// @Summary Get organization information
// @Description Get organization information
// @Tags Organization
// @Accept json
// @Produce json
// @Success 200 string {string} json "{"organizations": "organizations"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get organization"}"
// @Router /api/v1/organizations/ [get]
func GetOrganizationInformation(c *gin.Context) {
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

	var organization models.Organization
	result := DB.First(&organization, userOrganizationID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve organization"})
		return
	}

	var usersCount int64
	var projectsCount int64

	DB.Model(&models.User{}).Where("organization_id = ?", userOrganizationID).Count(&usersCount)
	DB.Model(&models.Project{}).Where("organization_id = ?", userOrganizationID).Count(&projectsCount)

	organizationResponseBody := gin.H{
		"organization": gin.H{
			"ID":                organization.ID,
			"CreatedAt":         organization.CreatedAt,
			"Name":              organization.Name,
			"AliasName":         organization.AliasName,
			"EstablishmentDate": organization.EstablishmentDate,
			"Description":       organization.Description,
			"UserCount":         usersCount,
			"ProjectCount":      projectsCount,
			"Address":           organization.Address,
		},
	}

	c.JSON(http.StatusOK, organizationResponseBody)
}

type organizationBody struct {
	Name              string `gorm:"type:varchar(100);not null"`
	AliasName         string `gorm:"type:varchar(100)"`
	EstablishmentDate time.Time
	Description       string `gorm:"type:text"`
}

// UpdateOrganizationInformation godoc
// @Summary Update organization information
// @Description Update organization information
// @Tags Organization
// @Accept json
// @Produce json
// @Param Organization body organizationBody true "Organization"
// @Success 200 string {string} json "{"organizations": "organizations"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get organization"}"
// @Router /api/v1/organizations/ [put]
func UpdateOrganizationInformation(c *gin.Context) {
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
	fmt.Println("debug")
	// Retrieve organization from the database using the user's organization ID
	var organization models.Organization
	result := DB.First(&organization, userOrganizationID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve organization"})
		return
	}

	// Bind JSON data to organizationBody struct
	var requestBody organizationBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update organization fields
	organization.Name = requestBody.Name
	organization.AliasName = requestBody.AliasName
	organization.EstablishmentDate = requestBody.EstablishmentDate
	organization.Description = requestBody.Description

	// Save the updated organization back to the database
	DB.Save(&organization)

	c.JSON(http.StatusOK, gin.H{"organization": organization})
}

// ListProjects godoc
// @Summary List projects
// @Description List projects
// @Tags Organization
// @Accept json
// @Produce json
// @Success 200 string {string} json "{"projects": "projects"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to list projects"}"
// @Router /api/v1/organizations/projects [get]
func ListProjects(c *gin.Context) {
	userInContext, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	user := userInContext.(models.User)
	if user.OrganizationID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization ID from user"})
		return
	}
	var projects []models.Project

	if user.IsOrganizationAdmin {
		DB.Where("organization_id = ?", user.OrganizationID).Find(&projects)
		c.JSON(http.StatusOK, gin.H{"projects": projects})
	} else {
		DB.Joins("JOIN user_project_roles ON projects.id = user_project_roles.project_id").Where("user_project_roles.user_id = ?", user.ID).Find(&projects)
		c.JSON(http.StatusOK, gin.H{"projects": projects})
	}
}

// CreateNewProject godoc
// @Summary Create new project
// @Description Create new project for organization
// @Tags Organization
// @Accept json
// @Produce json
// @Param Project body projectBody true "Project"
// @Success 200 string {string} json "{"project": "project"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to create project"}"
// @Router /api/v1/organizations/projects [post]
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

// DeleteProject godoc
// @Summary Delete project
// @Description Delete project
// @Tags Organization
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"message": "Project deleted"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to delete project"}"
// @Router /api/v1/organizations/projects/{project_id} [delete]
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
