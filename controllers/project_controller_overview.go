package controllers

import (
	"log"
	"net/http"
	"parameter-store-be/models"
	"parameter-store-be/modules/github"
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
	result := DB.
		Preload("Stages").
		Preload("Environments").
		First(&project, projectID)
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
	Name          string `binding:"required" json:"name" `
	StartAt       string `binding:"required" json:"start_at" `
	Description   string `binding:"required" json:"description" `
	CurrentSprint string `binding:"required" json:"current_sprint" `
	AutoUpdate    bool   `json:"auto_update" `
	RepoURL       string `binding:"required" json:"repo_url" `
	RepoApiToken  string `binding:"required" json:"repo_api_token" `
}

func (pb projectBody) Print() {
	log.Println("name: ", pb.Name)
	log.Println("start_at: ", pb.StartAt)
	log.Println("description: ", pb.Description)
	log.Println("current_sprint: ", pb.CurrentSprint)
	log.Println("auto_update: ", pb.AutoUpdate)
	log.Println("repo_url: ", pb.RepoURL)
	log.Println("repo_api_token: ", pb.RepoApiToken)
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
	requestBody.Print()
	// Retrieve project ID from the URL
	projectID := c.Param("project_id")

	// Retrieve project from the database using the project ID
	var project models.Project
	result := DB.First(&project, projectID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project"})
		return
	}
	// parse string to time
	startAt, err := time.Parse("2006-01-02", requestBody.StartAt)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date"})
		return
	}
	// Update project fields
	project.Name = requestBody.Name
	project.StartAt = startAt
	project.Description = requestBody.Description
	project.CurrentSprint = requestBody.CurrentSprint
	project.RepoURL = requestBody.RepoURL
	if err := github.ValidateGithubRepo(project.RepoURL, project.RepoApiToken); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
// @Param UserRoleProject body controllers.AddUserToProject.UserRoleBody true "UserRoleProject"
// @Success 200 string {string} json "{"message": "User added to project"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to add user to project"}"
// @Router /api/v1/projects/{project_id}/overview/add-user [post]
func AddUserToProject(c *gin.Context) {
	// Bind JSON data to UserRoleProject struct
	type UserRoleBody struct {
		Username string ` json:"username"`
		Role     string ` json:"role"`
	}
	var urb UserRoleBody
	if err := c.ShouldBindJSON(&urb); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve project ID from the URL
	projectID := c.Param("project_id")
	parsedProjectID, err := strconv.ParseUint(projectID, 10, 64)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	org_id, exist := c.Get("org_id")
	if !exist {
		log.Println("Failed to get organization ID from user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization ID from user"})
		return
	}

	var addedUser models.User
	// Retrieve user from the database using the username
	result := DB.Where("username = ? AND organization_id = ?", urb.Username, org_id).First(&addedUser)
	if result.RowsAffected == 0 {

		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	var role models.Role
	// Retrieve role from the database using the role name
	result = DB.Where("name = ?", urb.Role).First(&role)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found"})
		return
	}

	var urp models.UserRoleProject
	// Check if the user is already in the project
	result = DB.Where("user_id = ? AND project_id = ?", addedUser.ID, parsedProjectID).First(&urp)
	if result.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already in the project"})
		return
	}

	// Create a new user to project relationship
	urp = models.UserRoleProject{
		UserID:    addedUser.ID,
		ProjectID: uint(parsedProjectID),
		RoleID:    role.ID,
	}
	// Save the new user to project relationship to the database
	DB.Create(&urp)

	c.JSON(http.StatusOK, gin.H{"message": "User added to project"})
}

// RemoveUserFromProject godoc
// @Summary Remove user from project
// @Description Remove user from project
// @Tags Project Detail / Overview
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param user_id path string true "User ID"
// @Success 200 string {string} json "{"message": "User removed from project"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to remove user from project"}"
// @Router /api/v1/projects/{project_id}/overview/remove-user/{user_id} [delete]
func RemoveUserFromProject(c *gin.Context) {
	// Retrieve project ID and user ID from the URL
	projectID := c.Param("project_id")
	userID := c.Param("user_id")

	// Retrieve user to project relationship from the database using the project ID and user ID
	var urp models.UserRoleProject
	result := DB.Where("user_id = ? AND project_id = ?", userID, projectID).First(&urp)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is not in the project"})
		return
	}

	// Delete the user to project relationship from the database
	DB.Delete(&urp)

	c.JSON(http.StatusOK, gin.H{"message": "User removed from project"})
}
