package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/parameters [get]
func GetProjectParameters(c *gin.Context) {
	projectID := c.Param("project_id")
	page := c.Query("page")
	limit := c.Query("limit")
	stages := c.QueryArray("stages[]")
	environments := c.QueryArray("environments[]")
	version := c.Query("version")
	singleStage := c.Query("stages")
	singleEnvironment := c.Query("environments")
	filteredStage := append(stages, singleStage)
	filteredEnvironment := append(environments, singleEnvironment)
	// fmt.Println("Debug version", version)
	// Get project by ID
	var project models.Project
	var selectedVersion models.Version
	if version != "" {
		if err := DB.Preload("Versions", "number = ?", version).
			// where parameter is not archived
			Preload("Versions.Parameters", "is_archived = ?", false).
			Preload("Versions.Parameters.Stage").
			Preload("Versions.Parameters.Environment").
			First(&project, projectID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
			return
		}
		//debug version of project
		selectedVersion = project.Versions[0]
	} else {
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
		selectedVersion = project.LatestVersion
	}
	// Filter parameters by stages
	if len(filteredStage) > 0 {
		var filteredParameters []models.Parameter
		for _, stage := range filteredStage {
			for _, parameter := range selectedVersion.Parameters {
				if parameter.Stage.Name == stage {
					filteredParameters = append(filteredParameters, parameter)
				}
			}
		}
		selectedVersion.Parameters = filteredParameters
	}
	// Filter parameters by environments
	if len(filteredEnvironment) > 0 {
		var filteredParameters []models.Parameter
		for _, environment := range filteredEnvironment {
			for _, parameter := range selectedVersion.Parameters {
				if parameter.Environment.Name == environment {
					filteredParameters = append(filteredParameters, parameter)
				}
			}
		}
		selectedVersion.Parameters = filteredParameters
	}

	totalParam := len(selectedVersion.Parameters)
	var paginatedListParams []models.Parameter
	if page != "" && limit != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
			return
		}
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
			return
		}
		paginatedListParams = paginationDataParam(selectedVersion.Parameters, pageInt, limitInt)
	}
	c.JSON(http.StatusOK, gin.H{
		"parameters": paginatedListParams,
		"total":      totalParam,
	})
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
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
// @Failure 500 string {string} json "{"error": "Failed to get latest parameter"}"
// @Router /api/v1/projects/{project_id}/parameters/ [get]
func GetLatestParameters(c *gin.Context) {
	projectID, exist := c.Get("project_id")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project ID from user"})
		return
	}

	var project models.Project
	if err := DB.
		Preload("LatestVersion").
		Preload("LatestVersion.Parameters", "is_archived = ? ", false).
		Preload("LatestVersion.Parameters.Stage").
		Preload("LatestVersion.Parameters.Environment").
		First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}
	latestVersion := project.LatestVersion

	c.JSON(http.StatusOK, gin.H{"parameters": latestVersion.Parameters})
}

// Download lastest parameters in project, send file to client
// @Summary Download lastest parameters in project
// @Description Download lastest parameters in project
// @Tags Project Detail / Parameters
// @Accept json
// @Produce octet-stream
// @Param project_id path string true "Project ID"
// @Success 200 {array} models.Parameter
// @Security ApiKeyAuth
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get latest parameter"}"
// @Router /api/v1/projects/{project_id}/parameters/download [get]
func DownloadLatestParameters(c *gin.Context) {
	projectID := c.Param("project_id")
	queryStages := c.QueryArray("stages[]")
	queryEnvironments := c.QueryArray("environments[]")
	version := c.Query("version")
	singleStage := c.Query("stages")
	singleEnvironment := c.Query("environments")
	filteredStage := append(queryStages, singleStage)
	filteredEnvironment := append(queryEnvironments, singleEnvironment)

	var project models.Project
	var selectedVersion models.Version
	if version != "" {
		if err := DB.Preload("Versions", "number = ?", version).
			// where parameter is not archived
			Preload("Versions.Parameters", "is_archived = ?", false).
			Preload("Versions.Parameters.Stage", "is_archived = ?", false).
			Preload("Versions.Parameters.Environment", "is_archived = ?", false).
			Preload("Stages", "is_archived = ?", false).
			Preload("Environments", "is_archived = ?", false).
			First(&project, projectID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
			return
		}
		//debug version of project
		selectedVersion = project.Versions[0]
	} else {
		if err := DB.
			Preload("LatestVersion").
			// where parameter is not archived
			Preload("LatestVersion.Parameters", "is_archived = ?", false).
			Preload("LatestVersion.Parameters.Stage", "is_archived = ?", false).
			Preload("LatestVersion.Parameters.Environment", "is_archived = ?", false).
			Preload("Stages", "is_archived = ?", false).
			Preload("Environments", "is_archived = ?", false).
			First(&project, projectID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
			return
		}
		selectedVersion = project.LatestVersion
	}
	// Filter parameters by stages
	if len(filteredStage) > 0 {
		var filteredParameters []models.Parameter
		for _, stage := range filteredStage {
			for _, parameter := range selectedVersion.Parameters {
				if parameter.Stage.Name == stage {
					filteredParameters = append(filteredParameters, parameter)
				}
			}
		}
		selectedVersion.Parameters = filteredParameters
	}
	// Filter parameters by environments
	if len(filteredEnvironment) > 0 {
		var filteredParameters []models.Parameter
		for _, environment := range filteredEnvironment {
			for _, parameter := range selectedVersion.Parameters {
				if parameter.Environment.Name == environment {
					filteredParameters = append(filteredParameters, parameter)
				}
			}
		}
		selectedVersion.Parameters = filteredParameters
	}

	parameters := selectedVersion.Parameters

	// Create a new file
	filepath := fmt.Sprintf("parameters-%s-Ver.%s.txt", project.Name, selectedVersion.Number)
	file, err := os.Create(filepath) // format KEY=VALUE is paramter.Name=parameter.Value
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("######## Project: %s\n######## Version: %s \n", project.Name, selectedVersion.Number))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write project Name to file"})
		return
	}
	for _, environment := range project.Environments {
		if !isIn(filteredEnvironment, environment.Name) {
			continue
		}
		_, err := file.WriteString(fmt.Sprintf("\n########## ENVIRONMENT : %s ###########\n", environment.Name))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write environment to file"})
			return
		}
		for _, stage := range project.Stages {
			if !isIn(filteredStage, stage.Name) {
				continue
			}
			_, err := file.WriteString(fmt.Sprintf("\n#### STAGE : %s\n", stage.Name))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write stage to file"})
				return
			}
			for _, parameter := range parameters {
				if parameter.EnvironmentID == environment.ID && parameter.StageID == stage.ID {
					_, err = file.WriteString(fmt.Sprintf("%s=%s\n", parameter.Name, parameter.Value))
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
						return
					}
				}
			}
		}
		_, err = file.WriteString("\n##############################################\n")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write end file"})
			return
		}
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filepath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(filepath)
	err = os.Remove(filepath)
	if err != nil {
		log.Println("Failed to remove file")
	}
}

// check if value is in array
func isIn(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
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
// @Security ApiKeyAuth
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
		IsApplied:     false,
		Description:   newParameterBody.Description,
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

	// rerun github actions workflow if project.AutoUpdate is true
	if project.AutoUpdate {
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
	} else {
		projectLogByUser(newParameter.ProjectID, "Create Parameter", fmt.Sprint("Created parameter ", newParameter.Name), http.StatusCreated, 0, u.ID)
		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": "Parameter created",
		})
		return
	}
	// projectLogByUser(newParameter.ProjectID, "Create Parameter", fmt.Sprint("Created parameter ", newParameter.Name), responseStatusCode, latency, u.ID)
	// c.JSON(responseStatusCode, gin.H{
	// 	"status":  responseStatusCode,
	// 	"latency": latency,
	// 	"message": message})
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
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
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

	var project models.Project
	if err := DB.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}

	var parameter models.Parameter
	if err := DB.Where("project_id = ? AND id = ?", projectID, parameterID).First(&parameter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get parameter"})
		return
	}
	parameter.IsArchived = true
	parameter.ArchivedBy = u.Username
	parameter.ArchivedAt = time.Now()
	parameter.IsApplied = false
	if err := DB.Save(&parameter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to archive parameter"})
		return
	}
	// rerun github actions workflow
	if project.AutoUpdate {

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
	} else {
		projectLogByUser(parameter.ProjectID, "Archive Parameter", fmt.Sprint("Archived parameter ", parameter.Name), http.StatusCreated, 0, u.ID)
		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": "Parameter archived",
		})
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
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/parameters/{parameter_id}/unarchive [put]
func UnarchiveParameter(c *gin.Context) {
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	// Type assertion to extract user ID
	u := user.(models.User)
	projectID := c.Param("project_id")
	parameterID := c.Param("parameter_id")

	var project models.Project
	if err := DB.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}
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
	if project.AutoUpdate {

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
	} else {
		projectLogByUser(parameter.ProjectID, "Unarchive Parameter", fmt.Sprint("Unarchived parameter ", parameter.Name), http.StatusCreated, 0, u.ID)
		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": "Parameter unarchived",
		})
	}
	// projectLogByUser(parameter.ProjectID, "Unarchive Parameter", fmt.Sprint("Unarchived parameter ", parameter.Name), http.StatusCreated, latency, 0)
	// c.JSON(responseStatusCode, gin.H{
	// 	"status":  responseStatusCode,
	// 	"latency": latency,
	// 	"message": message})
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
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/parameters/{parameter_id} [put]
func UpdateParameter(c *gin.Context) {
	startTime := time.Now()
	projectID := c.Param("project_id")
	parameterID := c.Param("parameter_id")

	// get user from context
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	// modeling user
	u := user.(models.User)
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
	// fmt.Println("Debug updateParameterBody", updateParameterBody)

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
	parameter.IsApplied = false
	parameter.EditedAt = time.Now().UTC()
	if err := DB.Save(&parameter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update parameter"})
		return
	}
	// DB.Save(&parameter)

	// debug currentParameter and parameter
	if !project.AutoUpdate {
		latency := time.Since(startTime)
		projectLogByUser(parameter.ProjectID, "Update Parameter", fmt.Sprint("Updated parameter ", currentParameter.Name), http.StatusCreated, latency, u.ID)
		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": "Parameter updated",
		})
		return
	}
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
		projectLogByUser(parameter.ProjectID, "Update Parameter", fmt.Sprint("Updated parameter ", currentParameter.Name), http.StatusCreated, l, u.ID)
		c.JSON(responseStatusCode, gin.H{
			"status":  responseStatusCode,
			"latency": latency,
			"message": message})
		return
	}
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
		// Preload("Agents.Workflow").
		First(&project, updatedProjectID).Error; err != nil {
		return http.StatusInternalServerError, 0, "Failed to get project to rerun cicd", err
	}
	if len(project.Agents) == 0 {
		return http.StatusBadRequest, 0, "Failed to get agents to rerun CICD: no agents available", nil
	} else {
		usedAgent = project.Agents[0]
		// log.Println("Agents of project", project.Agents)
		if usedAgent.WorkflowName == "" {
			return http.StatusBadRequest, 0, "Failed to get agent workflow name to rerun cicd", nil
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
	// log.Println("lastAttemptNumber in rerunCICDWorkflow", lastAttemptNumber)
	// log.Println("lastWorkflowRunID in rerunCICDWorkflow", lastWorkflowRunID)
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

// ApplyParametersInProject godoc
// @Summary Apply parameters in project
// @Description Apply parameters in project
// @Tags Project Detail / Parameters
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"message": "Parameters applied"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to apply parameters"}"
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/parameters/apply [post]
func ApplyParametersInProject(c *gin.Context) {
	startTime := time.Now()
	projectID := c.Param("project_id")
	// parse project ID to uint
	projectIDUint64, err := strconv.ParseUint(projectID, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse project ID to uint"})
		return
	}
	projectIDUint := uint(projectIDUint64)

	var project models.Project
	if err := DB.
		Preload("LatestVersion").
		Preload("LatestVersion.Parameters").
		First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}
	latestVersion := project.LatestVersion
	// Get all parameters of the latest version
	parameters := latestVersion.Parameters
	// Apply all parameters
	// for _, parameter := range parameters {
	// 	parameter.IsApplied = true
	// }

	// find workflow ID of agent which is matched with stage and environment of un-IsApplied parameter
	var usedAgent models.Agent
	for _, parameter := range parameters {
		if !parameter.IsApplied {
			if err := DB.
				Preload("Workflow").
				Preload("Stage").
				Preload("Environment").
				Where("stage_id = ? AND environment_id = ?", parameter.StageID, parameter.EnvironmentID).
				First(&usedAgent).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get agent"})
				return
			}
			break
		}
	}
	responseStatusCode, latency, message, err := rerunCICDWorkflow(projectIDUint, usedAgent.StageID, usedAgent.EnvironmentID)
	latency = time.Since(startTime)
	if responseStatusCode == 403 {
		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"latency": latency,
			"message": "Parameters applied, but failed to rerun workflow: Workflow is already running. Check github actions of the project's repo.",
		})
		return
	}
	if err != nil {
		c.JSON(responseStatusCode, gin.H{"error": message})
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
	projectLogByUser(project.ID, "Apply Parameters", "Parameters applied", http.StatusCreated, latency, u.ID)
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"latency": latency,
		"message": "Parameters applied. Started rerun cicd. Check github actions of the project's repo.",
	})
}
