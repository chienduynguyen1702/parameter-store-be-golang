package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"parameter-store-be/models"
	"parameter-store-be/modules/github"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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
	IsUsingAtFile string `json:"is_using_at_file"`
}

// GetProjectParameters godoc
// @Summary Get project parameters
// @Description Get project parameters
// @Tags Project Detail / Parameters
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param page query string false "Page number"
// @Param limit query string false "Limit number"
// @Param stages query array false "Stages"
// @Param environments query array false "Environments"
// @Param version query string false "Version"
// @Param search query string false "Search"
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
	search := c.Query("search")
	filteredStage := append(stages, singleStage)
	filteredEnvironment := append(environments, singleEnvironment)
	// fmt.Println("Debug version", version)
	// Get project by ID
	var project models.Project
	if err := DB.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}
	var parameters []models.Parameter
	query := DB.Table("parameters").Select("parameters.*").Preload("Stage").Preload("Environment")
	if version != "" {
		query = query.
			Joins("LEFT JOIN version_parameters ON version_parameters.parameter_id = parameters.id").
			Joins("LEFT JOIN versions ON versions.id = version_parameters.version_id").
			Where("versions.number = ?", version)
	} else {
		query = query.
			Joins("LEFT JOIN version_parameters ON version_parameters.parameter_id = parameters.id").
			Joins("LEFT JOIN versions ON versions.id = version_parameters.version_id").
			Where("versions.id = ?", project.LatestVersionID)
	}
	if len(filteredStage) > 0 && filteredStage[0] != "" {
		// fmt.Println("Debug filteredStage", filteredStage, len(filteredStage))
		query = query.
			Joins("LEFT JOIN stages ON parameters.stage_id = stages.id").
			Where("stages.name IN (?)", filteredStage)
	}
	if len(filteredEnvironment) > 0 && filteredEnvironment[0] != "" {
		// fmt.Println("Debug filteredEnvironment", filteredEnvironment, len(filteredEnvironment))
		query = query.
			Joins("LEFT JOIN environments ON parameters.environment_id = environments.id").
			Where("environments.name IN (?)", filteredEnvironment)
	}
	query = query.Where("parameters.is_archived = ? AND parameters.project_id = ?", false, projectID).
		Order("parameters.environment_id DESC").
		Order("parameters.stage_id DESC")

	if search != "" {
		query = query.Where("parameters.name LIKE ? or parameters.value LIKE ? ", "%"+search+"%", "%"+search+"%")
	}
	query.Find(&parameters)
	// fmt.Println("Debug query parameters", parameters)
	// for _, p := range parameters {
	// 	fmt.Printf("%s=%s\n", p.Name, p.Value)
	// }
	totalParam := len(parameters)
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
		paginatedListParams = paginationDataParam(parameters, pageInt, limitInt)
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
	var project models.Project
	if err := DB.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}
	if project.RepoApiToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get repo api token"})
		return
	}
	if project.RepoURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get repo URL"})
		return
	}
	githubRepository, err := github.ParseRepoURL(project.RepoURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to parse repo URL"})
		return
	}
	// Get all files in the repo
	resultSearch := FindCodeAndFileContentInRepo(githubRepository.Owner, githubRepository.Name, project.RepoApiToken, parameter.Name)
	if resultSearch == "" {
		log.Println("Failed to find parameter", parameter.Name)
		c.JSON(http.StatusOK, gin.H{"is_using_at_file": resultSearch})
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
		IsUsingAtFile: resultSearch,
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
		_, err := file.WriteString(fmt.Sprintf("\n########## ENVS : %s ###########\n", environment.Name))
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

	githubRepository, err := github.ParseRepoURL(project.RepoURL)
	if err != nil {
		log.Println("Failed to parse repo URL in project", project.Name)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse repo URL"})
		return
	}
	resultSearching := FindCodeInRepo(githubRepository.Owner, githubRepository.Name, project.RepoApiToken, newParameterBody.Name)
	if resultSearching == "" {
		resultSearching = "null"
		return
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
		EditedAt:      time.Now().UTC(),
		Description:   newParameterBody.Description,
		IsUsingAtFile: resultSearching,
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

	githubRepository, err := github.ParseRepoURL(project.RepoURL)
	if err != nil {
		log.Println("Failed to parse repo URL in project", project.Name)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse repo URL"})
		return
	}
	resultSearching := FindCodeInRepo(githubRepository.Owner, githubRepository.Name, project.RepoApiToken, parameter.Name)
	if resultSearching == "" {
		resultSearching = "null"
		return
	}
	parameter.IsUsingAtFile = resultSearching
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
		Preload("Agents.Workflow").
		First(&project, updatedProjectID).Error; err != nil {
		return http.StatusInternalServerError, 0, "Failed to get project to rerun cicd", err
	}
	if len(project.Agents) == 0 {
		return http.StatusBadRequest, 0, "Failed to get agents to rerun CICD: no agents available", nil
	} else {
		usedAgent = project.Agents[0]
		// log.Println("Agents of project", project.Agents)
		if usedAgent.Workflow.Name == "" {
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
	// log.Println("usedAgent.WorkflowName in rerunCICDWorkflow", usedAgent.Workflow.Name)
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

// DownloadExecelTemplateParameters godoc
// @Summary Download execel template parameters
// @Description Download execel template parameters
// @Tags Project Detail / Parameters
// @Accept json
// @Produce octet-stream
// @Param project_id path string true "Project ID"
// @Success 200 {array} models.Parameter
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/parameters/download-template [get]
func DownloadExecelTemplateParameters(c *gin.Context) {
	project_id := c.Param("project_id")

	var project models.Project
	if err := DB.
		Preload("Stages", "is_archived = ?", false).
		Preload("Environments", "is_archived = ?", false).
		First(&project, project_id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}

	// Create a new file
	// log.Println("Creating file")
	stages := project.Stages
	envs := project.Environments
	// fmt.Println("stages", stages)
	// fmt.Println("envs", envs)

	newFile := excelize.NewFile()

	// file data to first sheet of created file
	newFile.SetSheetName("Sheet1", "Parameters")

	// set header of file
	// | Parameter Name | Value | Description | Stage | Environment |
	newFile.SetCellValue("Parameters", "A1", "Parameter Name")
	newFile.SetCellValue("Parameters", "B1", "Value")
	newFile.SetCellValue("Parameters", "C1", "Description")
	newFile.SetCellValue("Parameters", "D1", "Stage")
	newFile.SetCellValue("Parameters", "E1", "Environment")

	// create template parameters by envs and stages in project
	row := 2
	for _, env := range envs {
		for _, stage := range stages {
			newFile.SetCellValue("Parameters", "A"+strconv.Itoa(row), "KEY_NAME"+strconv.Itoa(row))
			newFile.SetCellValue("Parameters", "B"+strconv.Itoa(row), "value"+strconv.Itoa(row))
			newFile.SetCellValue("Parameters", "C"+strconv.Itoa(row), "description"+strconv.Itoa(row))
			newFile.SetCellValue("Parameters", "D"+strconv.Itoa(row), stage.Name)
			newFile.SetCellValue("Parameters", "E"+strconv.Itoa(row), env.Name)
			row++
			newFile.SetCellValue("Parameters", "A"+strconv.Itoa(row), "KEY_NAME"+strconv.Itoa(row))
			newFile.SetCellValue("Parameters", "B"+strconv.Itoa(row), "value"+strconv.Itoa(row))
			newFile.SetCellValue("Parameters", "C"+strconv.Itoa(row), "description"+strconv.Itoa(row))
			newFile.SetCellValue("Parameters", "D"+strconv.Itoa(row), stage.Name)
			newFile.SetCellValue("Parameters", "E"+strconv.Itoa(row), env.Name)
			row++
		}
	}
	// Set active sheet of the workbook.
	// newFile.SetActiveSheet(0)

	// Save the file
	filepath := fmt.Sprintf("Param-Template-%s.xlsx", project.Name)
	if err := newFile.SaveAs(filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}
	newFile.Close()
	// set header to sen excel file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filepath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	// send file
	c.File(filepath)
	err := os.Remove(filepath)
	if err != nil {
		log.Println("Failed to remove file")
	}
}

type UploadFileParamContent struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Stage       string `json:"stage"`
	Environment string `json:"environment"`

	StageID       uint
	EnvironmentID uint
}

// UploadParameters godoc
// @Summary Upload parameters
// @Description Upload parameters
// @Tags Project Detail / Parameters
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param file formData controllers.UploadFileParamContent true "File"
// @Success 200 string {string} json "{"message": "Parameters uploaded"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to upload parameters"}"
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/parameters/upload [post]
func UploadParameters(c *gin.Context) {
	startTime := time.Now()
	projectID := c.Param("project_id")
	file, err := c.FormFile("uploadFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
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
	// get project by ID
	var project models.Project
	if err := DB.
		Preload("Stages", "is_archived = ?", false).
		Preload("Environments", "is_archived = ?", false).
		First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}

	// Get stages and environments of project
	stages := project.Stages
	envs := project.Environments

	// Open the file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()
	// Read the file
	xlsx, err := excelize.OpenReader(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	// Get all the rows in the Sheet1.
	rows, err := xlsx.GetRows("Parameters")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rows"})
		return
	}
	// get row bind to []UploadFileParamContent
	var uploadFileParamContents []UploadFileParamContent
	for i, row := range rows {
		if i == 0 {
			continue
		}
		findingStageID := findStageID(stages, row[3])
		if findingStageID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to find stage name %s in row %d", row[3], i)})
			return
		}
		findingEnvID := findEnvironmentID(envs, row[4])
		if findingEnvID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to find environment name %s in row %d", row[4], i)})
			return
		}
		uploadFileParamContent := UploadFileParamContent{
			Name:          row[0],
			Value:         row[1],
			Description:   row[2],
			Stage:         row[3],
			Environment:   row[4],
			StageID:       findingStageID,
			EnvironmentID: findingEnvID,
		}
		uploadFileParamContents = append(uploadFileParamContents, uploadFileParamContent)
	}

	// Get the latest version of project
	var latestVersion models.Version
	if err := DB.
		Preload("Parameters").
		Preload("Parameters.Stage").
		Preload("Parameters.Environment").
		Last(&latestVersion, "project_id = ?", projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get latest version"})
		return
	}
	// Get all parameters of the latest version
	parameters := latestVersion.Parameters
	// Check if the parameter is already exist in the latest version, if not, insert it else overwrite it
	insertCount := 0
	overwriteCount := 0
	for _, uploadFileParamContent := range uploadFileParamContents {
		var isExist bool
		for _, parameter := range parameters {
			if parameter.Name == uploadFileParamContent.Name &&
				parameter.StageID == uploadFileParamContent.StageID &&
				parameter.EnvironmentID == uploadFileParamContent.EnvironmentID {
				//
				isExist = true
				parameter.Value = uploadFileParamContent.Value
				parameter.EditedAt = time.Now().UTC()
				parameter.Description = uploadFileParamContent.Description
				overwriteCount++
				break
			}
		}
		if isExist {
			continue
		}
		newParameter := models.Parameter{
			Name:          uploadFileParamContent.Name,
			Value:         uploadFileParamContent.Value,
			ProjectID:     project.ID,
			StageID:       uploadFileParamContent.StageID,
			EnvironmentID: uploadFileParamContent.EnvironmentID,
			IsApplied:     false,
			Description:   uploadFileParamContent.Description,
			EditedAt:      time.Now().UTC(),
		}
		// Append the new parameter to the latest version's Parameters slice
		latestVersion.Parameters = append(latestVersion.Parameters, newParameter)
		insertCount++
	}
	// Save the new parameter to the database
	if err := DB.Save(&latestVersion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save parameters"})
		return
	}

	// rerun github actions workflow if project.AutoUpdate is true
	if project.AutoUpdate {
		// Get project by ID to get agent workflow name
		wg := &sync.WaitGroup{}
		for _, ufpc := range uploadFileParamContents {
			wg.Add(1)
			go rerunCICDWorkflowWithGoRoutine(project.ID, ufpc.StageID, ufpc.EnvironmentID, u.ID, wg)
		}
		wg.Wait()
	} else {
		projectLogByUser(project.ID, "Upload Parameters", "Parameters uploaded", http.StatusCreated, 0, u.ID)
		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": "Parameters uploaded",
		})
		return
	}
	latency := time.Since(startTime)
	projectLogByUser(project.ID, "Upload Parameters", "Parameters uploaded", http.StatusCreated, latency, u.ID)
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"latency": latency,
		"message": "Parameters uploaded. Started rerun cicd. Check github actions of the project's repo.",
	})

}

func findStageID(stages []models.Stage, stageName string) uint {
	for _, stage := range stages {
		if stage.Name == stageName {
			return stage.ID
		}
	}
	return 0
}

func findEnvironmentID(envs []models.Environment, envName string) uint {
	for _, env := range envs {
		if env.Name == envName {
			return env.ID
		}
	}
	return 0
}

func rerunCICDWorkflowWithGoRoutine(projectID uint, stageID uint, environmentID uint, userID uint, wg *sync.WaitGroup) {
	// Get project by ID to get agent workflow name
	responseStatusCode, _, _, err := rerunCICDWorkflow(projectID, stageID, environmentID)
	if responseStatusCode == 403 {
		log.Println("Parameters uploaded, but failed to rerun workflow: Workflow is already running. Check github actions of the project's repo.")
		return
	}
	if err != nil {
		log.Println("Failed to rerun workflow")
		return
	}
	projectLogByUser(projectID, "Rerun CICD", "Parameters uploaded", http.StatusCreated, 0, userID)

	wg.Done()
}

type ParamPosition struct {
	ParameterName string
	Path          []struct {
		FileName   string
		LineNumber []int
	}
}

// SearchParameterInRepo godoc
// @Summary Search parameter in repo
// @Description Search parameter in repo
// @Tags Project Detail / Parameters
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param parameter_id path string true "Parameter ID"
// @Success 200 {array} models.Parameter
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to search parameters"}"
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/parameters/{parameter_id}/search-in-repo [get]
func SearchParameterInRepo(c *gin.Context) {
	projectID := c.Param("project_id")
	parameterID := c.Param("parameter_id")

	var project models.Project
	if err := DB.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}
	if project.RepoApiToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get repo api token"})
		return
	}
	if project.RepoURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get repo URL"})
		return
	}
	githubRepository, err := github.ParseRepoURL(project.RepoURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to parse repo URL"})
		return
	}
	// Get parameter by ID
	var parameter models.Parameter
	if err := DB.Where("project_id = ? AND id = ?", projectID, parameterID).First(&parameter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get parameter in project"})
		return
	}
	// Get parameter name
	parameterName := parameter.Name
	// Search parameter in repo
	searchResult, err := github.SearchCodeInRepo(githubRepository.Owner, githubRepository.Name, project.RepoApiToken, parameterName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search parameter in repo"})
		return
	}
	// var pathFiles []string
	var pP ParamPosition
	pP.ParameterName = parameterName

	// wait group
	var wg sync.WaitGroup
	for _, item := range searchResult.Items {
		wg.Add(1)
		go func(item github.SearchCodeInRepoItem) {
			defer wg.Done()
			// log.Println(item.Path)
			// pathFiles = append(pathFiles, item.Path)

			fileAsString, err := github.GetFileContent(githubRepository.Owner, githubRepository.Name, item.Path, project.RepoApiToken)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search parameter in repo"})
				return
			}
			lines := github.FindStringByLineNumber(fileAsString, parameterName)
			pP.Path = append(pP.Path, struct {
				FileName   string
				LineNumber []int
			}{FileName: item.Path, LineNumber: lines})

		}(item)

	}

	wg.Wait()
	// log.Println(pP)

	c.JSON(http.StatusOK, gin.H{"searchResult": searchResult})
}

// GetFileContent godoc
// @Summary Get file content
// @Description Get file content
// @Tags Project Detail / Parameters
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param parameter_id path string true "Parameter ID"
// @Param path query string false "Path file"
// @Success 200 string {string} json "{"file as string": "fileAsString"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get file content"}"
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/parameters/{parameter_id}/get-file-content [get]
func TestGetFileContent(c *gin.Context) {
	projectID := c.Param("project_id")
	parameterID := c.Param("parameter_id")
	path := c.Query("path")

	var project models.Project
	if err := DB.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}
	if project.RepoApiToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get repo api token"})
		return
	}
	if project.RepoURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get repo URL"})
		return
	}
	githubRepository, err := github.ParseRepoURL(project.RepoURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to parse repo URL"})
		return
	}
	// Get parameter by ID
	var parameter models.Parameter
	if err := DB.Where("project_id = ? AND id = ?", projectID, parameterID).First(&parameter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get parameter in project"})
		return
	}
	// Get parameter name
	// parameterName := parameter.Name
	// Search parameter in repo
	fileAsString, err := github.GetFileContent(githubRepository.Owner, githubRepository.Name, path, project.RepoApiToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search parameter in repo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"file as string": fileAsString})
}

type UsingAt struct {
	FileName     string `json:"file_name"`
	FileHTMLPath string `json:"file_html_path"`
	LineNumber   []int  `json:"line_number"`
	FileContent  string `json:"file_content"`
}
type CheckParamUsingBody struct {
	ParameterName string `json:"parameter_name"`
}

// GetFileContent godoc
// @Summary Get file content
// @Description Get file content
// @Tags Project Detail / Parameters
// @Accept json
// @Produce json
// @Param parameter_name body controllers.CheckParamUsingBody false "parameter name"
// @Success 200 string {string} json "{"file as string": "fileAsString"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get file content"}"
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/parameters/check-using [POST]
func CheckParameterUsing(c *gin.Context) {
	projectID := c.Param("project_id")

	var checkParamUsingBody CheckParamUsingBody
	if err := c.ShouldBindJSON(&checkParamUsingBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind json"})
		return
	}

	var project models.Project
	if err := DB.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}
	if project.RepoApiToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get repo api token"})
		return
	}
	if project.RepoURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get repo URL"})
		return
	}
	githubRepository, err := github.ParseRepoURL(project.RepoURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to parse repo URL"})
		return
	}
	// Get all files in the repo
	resultSearch := FindCodeAndFileContentInRepo(githubRepository.Owner, githubRepository.Name, project.RepoApiToken, checkParamUsingBody.ParameterName)
	if resultSearch == "" {
		log.Println("Failed to find parameter", checkParamUsingBody.ParameterName)
		c.JSON(http.StatusOK, gin.H{"is_using_at_file": resultSearch})
		return
	}
	// log.Println(checkParamUsingBody.ParameterName, "is in using ", resultSearch)
	c.JSON(http.StatusOK, gin.H{"is_using_at_file": resultSearch})
}

func AutoUpdateParameterUsingInFile() {
	for {
		log.Println("FindParameterUsingInFile")
		var projects []models.Project
		if err := DB.Find(&projects).Error; err != nil {
			log.Println("Failed to get projects")
			return
		}
		for _, project := range projects {
			fmt.Println("Starting Project", project.Name)

			if project.RepoApiToken == "" {
				log.Println("Failed to get repo api token in project", project.Name)
				continue
			}
			if project.RepoURL == "" {
				log.Println("Failed to get repo URL in project", project.Name)
				continue
			}
			githubRepository, err := github.ParseRepoURL(project.RepoURL)
			if err != nil {
				log.Println("Failed to parse repo URL in project", project.Name)
				continue
			}
			// Get all parameters of the project
			var parameters []models.Parameter
			if err := DB.Where("project_id = ?", project.ID).Find(&parameters).Error; err != nil {
				log.Println("Failed to get parameters in project", project.Name)
				continue
			}
			// Get all content of the files
			for _, param := range parameters {
				log.Println("")
				log.Println("Starting Parameter", param.Name)
				time.Sleep(6 * time.Second)

				// Get all files in the repo
				resultSearch := FindCodeInRepo(githubRepository.Owner, githubRepository.Name, project.RepoApiToken, param.Name)
				if resultSearch == "" {
					log.Println("Failed to find parameter", param.Name)
					continue
				}

				param.IsUsingAtFile = resultSearch
				if err := DB.Save(&param).Error; err != nil {
					log.Println("Failed to save parameter", param.Name)
					return
				}
				log.Println("Saved parameter", param.Name)
			}
			log.Println("=>>>>>>>>>> Finished project", project.Name)
		}
		log.Println("Finished")
	}
}

func FindCodeInRepo(owner, repo, token, paramName string) string {

	// Get all files in the repo
	searchCodeInRepoResponse, err := github.SearchCodeInRepo(owner, repo, token, paramName)
	if err != nil {
		log.Println("Failed to get files in repo")
		return ""
	}
	// total := len(searchCodeInRepoResponse.Items)
	// log.Println("Total files in repo", total)
	var usingAtTotal []UsingAt
	for _, item := range searchCodeInRepoResponse.Items {
		fileContent, err := github.GetFileContent(owner, repo, item.Path, token)
		if err != nil {
			log.Println("Failed to get file content")
			return ""
		}

		// Find parameter in file
		lines := github.FindStringByLineNumber(fileContent, paramName)
		if len(lines) > 0 {
			// log.Printf("Found parameter %s in file %s\n", paramName, item.Path)
		} else {
			continue
		}
		usingAt := UsingAt{
			FileName:     item.Path,
			FileHTMLPath: item.HTMLURL,
			LineNumber:   lines,
			// FileContent:  fileContent,
		}
		usingAtTotal = append(usingAtTotal, usingAt)
	}

	// remove redundunt usingAt
	// marshal slice of usingAt to json string and set to param.IsUsingAtFile
	usingAtJSON, err := json.Marshal(usingAtTotal)
	if err != nil {
		log.Println("Failed to marshal usingAt to json")
		return ""
	}
	// log.Println(paramName, "is in using ", string(usingAtJSON))
	return string(usingAtJSON)
}

func FindCodeAndFileContentInRepo(owner, repo, token, paramName string) string {

	// Get all files in the repo
	searchCodeInRepoResponse, err := github.SearchCodeInRepo(owner, repo, token, paramName)
	if err != nil {
		log.Println("Failed to get files in repo")
		return ""
	}
	// total := len(searchCodeInRepoResponse.Items)
	// log.Println("Total files in repo", total)
	var usingAtTotal []UsingAt
	for _, item := range searchCodeInRepoResponse.Items {
		fileContent, err := github.GetFileContent(owner, repo, item.Path, token)
		// log.Println("FileContent", fileContent)
		if err != nil {
			log.Println("Failed to get file content")
			return ""
		}

		// Find parameter in file
		lines := github.FindStringByLineNumber(fileContent, paramName)
		if len(lines) > 0 {
			// log.Printf("Found parameter %s in file %s\n", paramName, item.Path)
		} else {
			continue
		}
		usingAt := UsingAt{
			FileName:     item.Path,
			FileHTMLPath: item.HTMLURL,
			LineNumber:   lines,
			FileContent:  fileContent,
		}
		usingAtTotal = append(usingAtTotal, usingAt)
	}

	// remove redundunt usingAt
	// marshal slice of usingAt to json string and set to param.IsUsingAtFile
	usingAtJSON, err := json.Marshal(usingAtTotal)
	if err != nil {
		log.Println("Failed to marshal usingAt to json")
		return ""
	}
	// log.Println(paramName, "is in using ", string(usingAtJSON))
	return string(usingAtJSON)
}
