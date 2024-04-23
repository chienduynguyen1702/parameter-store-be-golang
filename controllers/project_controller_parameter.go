package controllers

import (
	"fmt"
	"log"
	"net/http"
	"parameter-store-be/models"
	"parameter-store-be/modules/github"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type parameterResponse struct {
	ID            uint   `json:"id"`
	StageID       uint   `json:"stage_id"`
	Stage         string `json:"stage"`
	EnvironmentID uint   `json:"environment_id"`
	Environment   string `json:"environment"`
	Name          string `json:"name"`
	Value         string `json:"value"`
	ProjectID     uint   `json:"project_id"`
	Description   string `json:"description"`
}

// GetProjectParameters godoc
// @Summary Get project parameters
// @Description Get project parameters
// @Tags Project Detail / Parameters
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {array} models.Parameter
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get project parameters"}"
// @Router /api/v1/projects/{project_id}/parameters [get]
func GetProjectParameters(c *gin.Context) {
	projectID := c.Param("project_id")

	// Get project by ID
	var project models.Project
	if err := DB.
		Preload("LatestVersion").
		// where parameter is not archived
		Preload("LatestVersion.Parameters", "is_archived = ?", false).
		Preload("LatestVersion.Parameters.Stage").
		Preload("LatestVersion.Parameters.Environment").
		First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"parameters": project.LatestVersion.Parameters})
}

// GetParameterByID godoc
// @Summary Get parameter by ID
// @Description Get parameter by ID
// @Tags Project Detail / Parameters
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param parameter_id path string true "Parameter ID"
// @Success 200 {object} models.Parameter
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get parameter"}"
// @Router /api/v1/projects/{project_id}/parameters/{parameter_id} [get]
func GetParameterByID(c *gin.Context) {
	projectID := c.Param("project_id")
	parameterID := c.Param("parameter_id")

	var parameter models.Parameter
	if err := DB.Preload("Stage").Preload("Environment").
		Where("project_id = ? AND id = ?", projectID, parameterID).First(&parameter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get parameter"})
		return
	}
	p := parameterResponse{
		ID:            parameter.ID,
		StageID:       parameter.StageID,
		Stage:         parameter.Stage.Name,
		EnvironmentID: parameter.EnvironmentID,
		Environment:   parameter.Environment.Name,
		Name:          parameter.Name,
		Value:         parameter.Value,
		ProjectID:     parameter.ProjectID,
		Description:   parameter.Description,
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"parameter": p,
		},
	})
}

// GetLatestParameters godoc
// @Summary Get latest parameter
// @Description Get latest parameter
// @Tags Project Detail / Parameters
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {array} models.Parameter
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get latest parameter"}"
// @Router /api/v1/projects/{project_id}/parameters/ [get]
func GetLatestParameters(c *gin.Context) {
	projectID, exist := c.Get("project_id")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project ID from user"})
		return
	}

	var project models.Project
	if err := DB.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}
	latestVersion := project.Versions[len(project.Versions)-1]

	c.JSON(http.StatusOK, gin.H{"parameters": latestVersion.Parameters})
}

// CreateParameter godoc
// @Summary Create new parameter
// @Description Create new parameter
// @Tags Project Detail / Parameters
// @Accept json
// @Produce json
// @Param Parameter body controllers.CreateParameter.createParameterRequestBody true "Parameter"
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"message": "Parameter created"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to create parameter"}"
// @Router /api/v1/projects/{project_id}/parameters [post]
func CreateParameter(c *gin.Context) {
	projectID := c.Param("project_id")

	// get user from context
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	// modeling user
	u := user.(models.User)
	type createParameterRequestBody struct {
		Name        string `json:"name" binding:"required"`
		Value       string `json:"value" binding:"required"`
		Stage       string `json:"stage"`
		Environment string `json:"environment"`
		Description string `json:"description"`
	}
	newParameterBody := createParameterRequestBody{}
	if err := c.ShouldBindJSON(&newParameterBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// get latest version of project
	var project models.Project
	if err := DB.
		Preload("Versions").
		Preload("Stages").
		Preload("Environments").
		Preload("Workflows").
		First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}
	latestVersion := project.Versions[len(project.Versions)-1]

	// var stage models.Stage
	// if err := DB.Where("name = ?", newParameterBody.Stage).First(&stage).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stage"})
	// 	return
	// }
	// var environment models.Environment
	// if err := DB.Where("name = ?", newParameterBody.Environment).First(&environment).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get environment"})
	// 	return
	// }
	stages := project.Stages
	var findingStage models.Stage
	for _, stage := range stages {
		if stage.Name == newParameterBody.Stage {
			findingStage = stage
			break
		}
	}
	env := project.Environments
	var findingEnvironment models.Environment
	for _, e := range env {
		if e.Name == newParameterBody.Environment {
			findingEnvironment = e
			break
		}
	}

	newParameter := models.Parameter{
		Name:          newParameterBody.Name,
		Value:         newParameterBody.Value,
		ProjectID:     project.ID,
		StageID:       findingStage.ID,
		EnvironmentID: findingEnvironment.ID,
		Stage:         findingStage,
		Environment:   findingEnvironment,
	}

	// Append the new parameter to the latest version's Parameters slice
	latestVersion.Parameters = append(latestVersion.Parameters, newParameter)
	// Save the new parameter to the database
	// if err := DB.Create(&newParameter).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create parameter"})
	// 	return
	// }
	if err := DB.Save(&latestVersion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update latest version"})
		return
	}

	// rerun github actions workflow
	// Get project by ID to get agent workflow name
	responseStatusCode, latency, message, err := rerunCICDWorkflow(newParameter.ProjectID, newParameter.StageID, newParameter.EnvironmentID)
	if responseStatusCode == 403 {
		responseStatusCode = http.StatusCreated
		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"latency": latency,
			"message": "Parameter updated, but failed to rerun workflow: Workflow is already running. Check github actions of the project's repo.",
		})
		return
	}
	if err != nil {
		projectLogByUser(newParameter.ProjectID, "Create Parameter", "Failed to create parameter", responseStatusCode, latency, u.ID)
		c.JSON(responseStatusCode, gin.H{"error": message})
		return
	}
	projectLogByUser(newParameter.ProjectID, "Create Parameter", fmt.Sprint("Created parameter ", newParameter.Name), responseStatusCode, latency, u.ID)
	c.JSON(responseStatusCode, gin.H{
		"status":  responseStatusCode,
		"latency": latency,
		"message": message})
}

// GetArchivedParameters godoc
// @Summary Get archived parameters
// @Description Get archived parameters
// @Tags Project Detail / Parameters
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {array} models.Parameter
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get archived parameters"}"
// @Router /api/v1/projects/{project_id}/parameters/archived [get]
func GetArchivedParameters(c *gin.Context) {
	projectID := c.Param("project_id")

	var parameters []models.Parameter
	DB.Preload("Stage").Preload("Environment").
		Where("project_id = ? AND is_archived = ?", projectID, true).Find(&parameters)
	c.JSON(http.StatusOK, gin.H{"parameters": parameters})
}

// ArchiveParameter godoc
// @Summary Archive parameter
// @Description Archive parameter
// @Tags Project Detail / Parameters
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param parameter_id path string true "Parameter ID"
// @Success 200 string {string} json "{"message": "Parameter archived"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to archive parameter"}"
// @Router /api/v1/projects/{project_id}/parameters/{parameter_id}/archive [put]
func ArchiveParameter(c *gin.Context) {
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	// Type assertion to extract user ID
	u := user.(models.User)

	projectID := c.Param("project_id")
	parameterID := c.Param("parameter_id")

	var parameter models.Parameter
	if err := DB.Where("project_id = ? AND id = ?", projectID, parameterID).First(&parameter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get parameter"})
		return
	}
	parameter.IsArchived = true
	parameter.ArchivedBy = u.Username
	parameter.ArchivedAt = time.Now()
	if err := DB.Save(&parameter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to archive parameter"})
		return
	}
	// rerun github actions workflow
	// Get project by ID to get agent workflow name
	responseStatusCode, latency, message, err := rerunCICDWorkflow(parameter.ProjectID, parameter.StageID, parameter.EnvironmentID)
	projectLogByUser(parameter.ProjectID, "Archive Parameter", fmt.Sprint("Archived parameter ", parameter.Name), http.StatusCreated, latency, u.ID)
	if responseStatusCode == 403 {
		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"latency": latency,
			"message": "Parameter updated, but failed to rerun workflow: Workflow is already running. Check github actions of the project's repo.",
		})
		return
	}
	if err != nil {
		c.JSON(responseStatusCode, gin.H{"error": message})
		return
	}
}

// UnarchiveParameter godoc
// @Summary Unarchive parameter
// @Description Unarchive parameter
// @Tags Project Detail / Parameters
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param parameter_id path string true "Parameter ID"
// @Success 200 string {string} json "{"message": "Parameter unarchived"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to unarchive parameter"}"
// @Router /api/v1/projects/{project_id}/parameters/{parameter_id}/unarchive [put]
func UnarchiveParameter(c *gin.Context) {
	projectID := c.Param("project_id")
	parameterID := c.Param("parameter_id")

	var parameter models.Parameter
	if err := DB.Where("project_id = ? AND id = ?", projectID, parameterID).First(&parameter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get parameter"})
		return
	}
	parameter.IsArchived = false
	parameter.ArchivedBy = ""
	parameter.ArchivedAt = time.Time{}
	if err := DB.Save(&parameter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unarchive parameter"})
		return
	}
	// rerun github actions workflow
	// Get project by ID to get agent workflow name
	responseStatusCode, latency, message, err := rerunCICDWorkflow(parameter.ProjectID, parameter.StageID, parameter.EnvironmentID)
	if responseStatusCode == 403 {
		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"latency": latency,
			"message": "Parameter updated, but failed to rerun workflow: Workflow is already running. Check github actions of the project's repo.",
		})
	}
	if err != nil {
		c.JSON(responseStatusCode, gin.H{"error": message})
		return
	}
	projectLogByUser(parameter.ProjectID, "Unarchive Parameter", fmt.Sprint("Unarchived parameter ", parameter.Name), http.StatusCreated, latency, 0)
	c.JSON(responseStatusCode, gin.H{
		"status":  responseStatusCode,
		"latency": latency,
		"message": message})
}

// UpdateParameter godoc
// @Summary Update parameter
// @Description Update parameter
// @Tags Project Detail / Parameters
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param parameter_id path string true "Parameter ID"
// @Param Parameter body controllers.UpdateParameter.updateParameterRequestBody true "Parameter"
// @Success 200 string {string} json "{"message": "Parameter updated"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to update parameter"}"
// @Router /api/v1/projects/{project_id}/parameters/{parameter_id} [put]
func UpdateParameter(c *gin.Context) {
	projectID := c.Param("project_id")
	parameterID := c.Param("parameter_id")

	type updateParameterRequestBody struct {
		Name        string `json:"name"`
		Value       string `json:"value"`
		Stage       string `json:"stage"`
		Environment string `json:"environment"`
		Description string `json:"description"`
	}
	updateParameterBody := updateParameterRequestBody{}
	if err := c.ShouldBindJSON(&updateParameterBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var parameter models.Parameter
	if err := DB.Where("project_id = ? AND id = ?", projectID, parameterID).First(&parameter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get parameter"})
		return
	}
	// preload stages and environments in project
	var project models.Project

	if err := DB.
		Preload("Stages").
		Preload("Environments").
		First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}
	// get stages and environments from project
	stages := project.Stages
	var findingStage models.Stage
	for _, stage := range stages {
		if stage.Name == updateParameterBody.Stage {
			findingStage = stage
			break
		}
	}

	env := project.Environments
	var findingEnvironment models.Environment
	for _, e := range env {
		if e.Name == updateParameterBody.Environment {
			findingEnvironment = e
			break
		}
	}

	//duplicate parameter to check if parameter is updated at Name or Value or Stage or Environment
	currentParameter := parameter
	if updateParameterBody.Name != "" {
		parameter.Name = updateParameterBody.Name
	}
	if updateParameterBody.Value != "" {
		parameter.Value = updateParameterBody.Value
	}
	if updateParameterBody.Description != "" {
		parameter.Description = updateParameterBody.Description
	}
	if updateParameterBody.Stage != "" {
		parameter.StageID = findingStage.ID
	}
	if updateParameterBody.Environment != "" {
		parameter.EnvironmentID = findingEnvironment.ID
	}

	if err := DB.Save(&parameter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update parameter"})
		return
	}
	// debug currentParameter and parameter

	// If parameter is updated at Name or Value or Stage or Environment
	// then rerun github actions workflow
	var l time.Duration
	if currentParameter.Name != parameter.Name ||
		currentParameter.Value != parameter.Value ||
		currentParameter.StageID != parameter.StageID ||
		currentParameter.EnvironmentID != parameter.EnvironmentID {
		// Get project by ID to get agent workflow name
		responseStatusCode, latency, message, err := rerunCICDWorkflow(parameter.ProjectID, parameter.StageID, parameter.EnvironmentID)

		l = latency
		if responseStatusCode == 403 {
			c.JSON(http.StatusCreated, gin.H{
				"status":  http.StatusCreated,
				"latency": latency,
				"message": "Parameter updated, but failed to rerun workflow: Workflow is already running. Check github actions of the project's repo.",
			})
		}
		if err != nil {
			c.JSON(responseStatusCode, gin.H{"error": message})
			return
		}
		c.JSON(responseStatusCode, gin.H{
			"status":  responseStatusCode,
			"latency": latency,
			"message": message})
		return
	}
	// get user from context
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	// modeling user
	u := user.(models.User)
	projectLogByUser(parameter.ProjectID, "Update Parameter", fmt.Sprint("Updated parameter ", currentParameter.Name), http.StatusCreated, l, u.ID)
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"latency": l,
		"message": "Parameter updated. Started rerun cicd. Check github actions of the project's repo.",
	})
}
func rerunCICDWorkflow(updatedProjectID uint, updatedStageID uint, updatedEnvironmentID uint) (int, time.Duration, string, error) {
	//debug
	// log.Println("rerunCICDWorkflow\n", "updatedProjectID", updatedProjectID, "updatedStageID", updatedStageID, "updatedEnvironmentID", updatedEnvironmentID)
	var project models.Project
	var usedAgent models.Agent
	if err := DB.
		Preload("Agents", "stage_id = ? AND environment_id = ?", updatedStageID, updatedEnvironmentID).
		Preload("Agents.Workflow").
		First(&project, updatedProjectID).Error; err != nil {
		return http.StatusInternalServerError, 0, "Failed to get project to rerun cicd", err
	}
	if project.Agents == nil {
		return http.StatusInternalServerError, 0, "Failed to get agent to rerun cicd", nil
	} else {
		usedAgent = project.Agents[0]
		// log.Println("Agents of project", project.Agents)
		if usedAgent.WorkflowName == "" {
			return http.StatusInternalServerError, 0, "Failed to get agent workflow name to rerun cicd", nil
		}
	}
	if project.RepoApiToken == "" {
		return http.StatusInternalServerError, 0, "Failed to get repo api token to rerun cicd", nil
	}
	if project.RepoURL == "" {
		return http.StatusInternalServerError, 0, "Failed to get repo URL to rerun cicd", nil
	}

	githubRepository, err := github.ParseRepoURL(project.RepoURL)
	if err != nil {
		return http.StatusNotFound, 0, "Failed to parse repo URL to rerun cicd", err
	}
	startTime := time.Now()
	responseStatusCode, responseMessage, err := github.RerunWorkFlow(githubRepository.Owner, githubRepository.Name, usedAgent.WorkflowName, project.RepoApiToken)
	latency := time.Since(startTime)

	lastWorkflowRunID, _, lastAttemptNumber, errAttempt := github.GetLastAttemptNumberOfWorkflowRun(githubRepository.Owner, githubRepository.Name, project.RepoApiToken, usedAgent.WorkflowName)
	if errAttempt != nil {
		log.Println("Failed to get last attempt number of workflow run")
	}
	log.Println("lastAttemptNumber in rerunCICDWorkflow", lastAttemptNumber)
	log.Println("lastWorkflowRunID in rerunCICDWorkflow", lastWorkflowRunID)
	//parse lastWorkflowRunID to uint
	lastWorkflowRunIDUint64, errParse := strconv.ParseUint(lastWorkflowRunID, 10, 64)
	if errParse != nil {
		log.Println("Failed to parse lastWorkflowRunID to uint")
	}
	//save workflow log
	workflowLog(usedAgent.WorkflowID, uint(lastWorkflowRunIDUint64), lastAttemptNumber)
	log.Println(responseMessage)
	if responseStatusCode == 403 {
		return 201, latency, fmt.Sprintf("Parameter updated. Failed to rerun workflow: Workflow is already running. Check github actions at %s/actions", project.RepoURL), nil
	}
	if err != nil {
		return http.StatusInternalServerError, 0, err.Error(), nil
	}
	return http.StatusCreated, latency, fmt.Sprintf("Parameter updated. Started rerun cicd. Check github actions at %s/actions", project.RepoURL), nil
}

// func rerunCICDWorkflowHandler(c *gin.Context) {
// 	var project models.Project
// 	var usedAgent models.Agent
// 	if err := DB.
// 		Preload("Agents", "stage_id = ? AND environment_id = ?", updatedStageID, updatedEnvironmentID).
// 		First(&project, updatedProjectID).Error; err != nil {
// 		return http.StatusInternalServerError, 0, "Failed to get project to rerun cicd", err
// 	}
// 	if project.Agents == nil {
// 		return http.StatusInternalServerError, 0, "Failed to get agent to rerun cicd", nil
// 	}
// 	if len(project.Agents) != 1 {
// 		return http.StatusInternalServerError, 0, "Failed to get agent to rerun cicd", nil
// 	} else {
// 		usedAgent = project.Agents[0]
// 		if usedAgent.WorkflowName == "" {
// 			return http.StatusInternalServerError, 0, "Failed to get agent workflow name to rerun cicd", nil
// 		}
// 	}
// 	if project.RepoApiToken == "" {
// 		return http.StatusInternalServerError, 0, "Failed to get repo api token to rerun cicd", nil
// 	}
// 	if project.RepoURL == "" {
// 		return http.StatusInternalServerError, 0, "Failed to get repo URL to rerun cicd", nil
// 	}

// 	githubRepository, err := github.ParseRepoURL(project.RepoURL)
// 	if err != nil {
// 		return http.StatusNotFound, 0, "Failed to parse repo URL to rerun cicd", err
// 	}
// 	startTime := time.Now()
// 	responseStatusCode, err := github.RerunWorkFlow(githubRepository.Owner, githubRepository.Name, usedAgent.WorkflowName, project.RepoApiToken)
// 	latency := time.Since(startTime)
// 	if responseStatusCode == 403 {
// 		return 403, latency, fmt.Sprintf("Parameter updated. Failed to rerun workflow: Workflow is already running. Check github actions at %s/actions", project.RepoURL), nil
// 	}
// 	if err != nil {
// 		return http.StatusInternalServerError, 0, err.Error(), nil
// 	}
// 	return http.StatusCreated, latency, fmt.Sprintf("Parameter updated. Started rerun cicd. Check github actions at %s/actions", project.RepoURL), nil
// }
