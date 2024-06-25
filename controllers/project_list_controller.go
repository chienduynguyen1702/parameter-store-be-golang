package controllers

import (
	"fmt"
	"log"
	"net/http"
	"parameter-store-be/models"
	"parameter-store-be/modules/github"
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
		Status    string `json:"status"`
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
			Status:    projects[i].Status,
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
		OrganizationID:  userOrganizationID,
		Name:            requestBody.Name,
		Description:     requestBody.Description,
		StartAt:         time.Now(),
		Status:          "In Progress",
		CurrentSprint:   "1",
		RepoURL:         "github.com/OWNER/REPO",
		LatestVersionID: 1,
	}
	// Save the new project to the database
	DB.Create(&project)
	initVersion := models.Version{
		Number:      "1.0.0",
		Name:        "1.0.0",
		Description: "Initial version",
	}
	DB.Create(&initVersion)

	project.LatestVersionID = initVersion.ID
	DB.Save(&project)
	// remove versionid 1 from project
	// DB.Delete()
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

// ListGithubRepos godoc
// @Summary List github repos
// @Description List github repos
// @Tags Project List
// @Accept json
// @Produce json
// @Success 200 string {string} json "{"total": "total", "repositories": "repositories"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to list repository by github user"}"
// @Security ApiKeyAuth
// @Router /api/v1/project-list/github-repos [get]
func ListGithubRepos(c *gin.Context) {
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
	repositories, err := github.GetUserRepos(user.Username, user.GithubAccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list repository by github user"})
		return
	}
	// for i := 0; i < len(repositories); i++ {
	// 	fmt.Println(repositories[i].Name)
	// }
	c.JSON(http.StatusOK, gin.H{
		"total":        len(repositories),
		"repositories": repositories,
	})
}

type ImportReposToProjectRequest struct {
	ID       uint   `json:"id"`
	Name     string `json:"name" `
	FullName string `json:"full_name" `
	HTMLURL  string `json:"html_url" `
}

// ImportReposToProject godoc
// @Summary Import repos to project
// @Description Import repos to project
// @Tags Project List
// @Accept json
// @Produce json
// @Param request body ImportReposToProjectRequest true "Import Repos To Project Request"
func ImportReposToProject(c *gin.Context) {
	var request []ImportReposToProjectRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	accessToken := user.(models.User).GithubAccessToken
	for i := 0; i < len(request); i++ {
		// fmt.Println(request[i].ID)
		// fmt.Println(request[i].Name)
		// fmt.Println(request[i].FullName)
		// fmt.Println(request[i].HTMLURL)
		githuburl := fmt.Sprintf("github.com/%s", request[i].FullName)
		err := CreateRepoToProject(userOrganizationID, request[i].Name, githuburl, accessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import repos to project"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Repos imported to project"})
}

func CreateRepoToProject(orgID uint, projectName, projectURL, projectAPIToken string) error {
	var findingProject models.Project
	DB.Where("repo_url = ?", projectURL).First(&findingProject)
	if findingProject.ID != 0 {
		fmt.Println("Project already exists")
		return nil
	}
	// Create a new project
	project := models.Project{
		OrganizationID: orgID,
		Name:           projectName,
		RepoURL:        projectURL,
		RepoApiToken:   projectAPIToken,
		Status:         "In Progress",
	}
	DB.Create(&project)
	newVersion := models.Version{
		Number:      "1.0.0",
		Name:        "1.0.0",
		Description: "Initial version",
		ProjectID:   project.ID,
	}
	DB.Create(&newVersion)
	// save latest version id to project
	project.LatestVersionID = newVersion.ID
	DB.Save(&project)

	// Create new stages
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
	repo, err := github.ParseRepoURL(projectURL)
	if err != nil {
		fmt.Println(err)
	}
	log.Println(repo.Owner, repo.Name, projectAPIToken)
	colaborator, err := github.ListRepositoryColaborator(repo.Owner, repo.Name, projectAPIToken)
	if err != nil {
		fmt.Println(err)
	}
	var urps []models.UserRoleProject
	for _, colab := range colaborator {
		fmt.Println(colab.Login, colab.Email, colab.Permissions.Admin)
		user := models.User{
			Username:       colab.Login,
			Email:          colab.Email,
			OrganizationID: orgID,
		}
		var urp models.UserRoleProject
		if err := DB.Where("email = ?", colab.Email).First(&user).Error; err != nil {
			DB.Create(&user)
		}

		if colab.Permissions.Admin { // check if user is admin in the repo
			urp = models.UserRoleProject{
				UserID:    user.ID,
				ProjectID: project.ID,
				RoleID:    2,
			}
		} else {
			urp = models.UserRoleProject{ // user is not admin in the repo
				UserID:    user.ID,
				ProjectID: project.ID,
				RoleID:    3,
			}
		}
		urps = append(urps, urp)
	}
	DB.Create(&urps)
	return nil
}
