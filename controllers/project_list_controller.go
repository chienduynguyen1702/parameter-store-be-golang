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
// @Tags Project List
// @Accept json
// @Produce json
// @Success 200 string {string} json "{"projects": "projects"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Security ApiKeyAuth
// @Failure 500 string {string} json "{"error": "Failed to list projects"}"
// @Router /api/v1/project-list/ [get]
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
		DB.Where("organization_id = ? AND is_archived != ? ", user.OrganizationID, true).Find(&projects)
	} else {
		DB.Joins("JOIN user_role_projects ON projects.id = user_role_projects.project_id").Where("user_role_projects.user_id = ? AND projects.is_archived != ? ", user.ID, true).Find(&projects)
	}

	type projectListResponse struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		UserCount int64  `json:"users_count"`
	}
	var projectListResponses []projectListResponse

	//count user in project
	for i := 0; i < len(projects); i++ {
		var userCount int64
		DB.Model(&models.UserRoleProject{}).Where("project_id = ?", projects[i].ID).Count(&userCount)
		projectListResponses = append(projectListResponses, projectListResponse{
			ID:        projects[i].ID,
			Name:      projects[i].Name,
			UserCount: userCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{"projects": projectListResponses})
}

// CreateNewProject godoc
// @Summary Create new project
// @Description Create new project for organization
// @Tags Project List
// @Accept json
// @Produce json
// @Param Project body projectBody true "Project"
// @Success 200 string {string} json "{"project": "project"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to create project"}"
// @Security ApiKeyAuth
// @Router /api/v1/project-list/ [post]
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

	// Create a new project3
	project := models.Project{
		OrganizationID: userOrganizationID,
		Name:           requestBody.Name,
		Description:    requestBody.Description,
		StartAt:        time.Now(),
		Status:         "In Progress",
		CurrentSprint:  "1",
		RepoURL:        "github.com/OWNER/REPO",
	}
	// Save the new project to the database
	DB.Create(&project)

	newStages := []models.Stage{
		{
			Name:        "Build",
			Description: "Build stage",
			ProjectID:   project.ID,
		},
		{
			Name:        "Test",
			Description: "Test stage",
			ProjectID:   project.ID,
		},
		{
			Name:        "Release",
			Description: "Release stage",
			ProjectID:   project.ID,
		},
		{
			Name:        "Deploy",
			Description: "Deploy stage",
			ProjectID:   project.ID,
		},
	}
	for _, stage := range newStages {
		DB.Create(&stage)
	}

	newEnvironment := []models.Environment{
		{

			Name:        "Development",
			Description: "Development environment",
			ProjectID:   project.ID,
		},
		{
			Name:        "Staging",
			Description: "Staging environment",
			ProjectID:   project.ID,
		},
		{
			Name:        "Production",
			Description: "Production environment",
			ProjectID:   project.ID,
		},
	}
	for _, environment := range newEnvironment {
		DB.Create(&environment)
	}
	newVersion := models.Version{
		Number:      "1.0.0",
		Name:        "1.0.0",
		ProjectID:   project.ID,
		Description: "Initial version",
	}
	DB.Create(&newVersion)

	c.JSON(http.StatusOK, gin.H{"project": project})
}

// DeleteProject godoc
// @Summary Delete project
// @Description Delete project
// @Tags Project List
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"message": "Project deleted"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to delete project"}"
// @Security ApiKeyAuth
// @Router /api/v1/project-list/{project_id} [delete]
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

// ListArchivedProjects godoc
// @Summary List archived projects
// @Description List archived projects
// @Tags Project List
// @Accept json
// @Produce json
// @Success 200 string {string} json "{"projects": "projects"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to list archived projects"}"
// @Security ApiKeyAuth
// @Router /api/v1/project-list/archived [get]
func ListArchivedProjects(c *gin.Context) {
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
		DB.Where("organization_id = ? AND is_archived = ?", user.OrganizationID, true).Find(&projects)
	} else {
		DB.Joins("JOIN user_role_projects ON projects.id = user_role_projects.project_id").Where("user_role_projects.user_id = ? AND projects.is_archived = ?", user.ID, true).Find(&projects)
	}

	type archivedProjectListResponse struct {
		ID         uint      `json:"id"`
		Name       string    `json:"name"`
		UserCount  int64     `json:"users_count"`
		ArchivedAt time.Time `json:"archived_at"`
		ArchivedBy string    `json:"archived_by"`
	}
	var projectListResponses []archivedProjectListResponse

	//count user in project
	for i := 0; i < len(projects); i++ {
		var userCount int64
		DB.Model(&models.UserRoleProject{}).Where("project_id = ?", projects[i].ID).Count(&userCount)
		projectListResponses = append(projectListResponses, archivedProjectListResponse{
			ID:         projects[i].ID,
			Name:       projects[i].Name,
			UserCount:  userCount,
			ArchivedAt: projects[i].ArchivedAt,
			ArchivedBy: projects[i].ArchivedBy,
		})
	}

	c.JSON(http.StatusOK, gin.H{"projects": projectListResponses})
}

// ArchiveProject godoc
// @Summary Archive project
// @Description Archive project
// @Tags Project List
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"message": "Project archived"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to archive project"}"
// @Security ApiKeyAuth
// @Router /api/v1/project-list/{project_id}/archive [put]
func ArchiveProject(c *gin.Context) {
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
	if project.IsArchived {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project is already archived"})
		return
	}

	// Archive the project
	project.Status = "Archived"
	project.IsArchived = true
	project.ArchivedAt = time.Now()
	project.ArchivedBy = user.(models.User).Email
	DB.Save(&project)

	c.JSON(http.StatusOK, gin.H{"message": "Project archived"})
}

// UnarchiveProject godoc
// @Summary Unarchive project
// @Description Unarchive project
// @Tags Project List
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"message": "Project unarchived"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to unarchive project"}"
// @Security ApiKeyAuth
// @Router /api/v1/project-list/{project_id}/unarchive [put]
func UnarchiveProject(c *gin.Context) {
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
	if !project.IsArchived {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project is not archived"})
		return
	}
	// unarchive the project
	project.Status = "In Progress"
	project.IsArchived = false
	project.ArchivedAt = time.Time{}
	project.ArchivedBy = ""
	DB.Save(&project)

	c.JSON(http.StatusOK, gin.H{"message": "Project unarchived"})
}
