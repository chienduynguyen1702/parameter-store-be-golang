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

type agentResponse struct {
	ID            uint `json:"id"`
	ProjectID     uint
	Name          string `gorm:"type:varchar(100);not null" json:"name"`
	StageID       uint   `gorm:"foreignKey:StageID;not null" json:"stage_id"`
	Stage         models.Stage
	EnvironmentID uint `gorm:"foreignKey:EnvironmentID;not null" json:"environment_id"`
	Environment   models.Environment
	WorkflowName  string `gorm:"type:varchar(100);not null" json:"workflow_name"`
}

// GetProjectAgents godoc
// @Summary Get agents of project
// @Description Get agents of project
// @Tags Project Detail / Agents
// @Accept json
// @Produce json
// @Param project_id path int true "Project ID"
// @Success 200 {array}	 controllers.agentResponse
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get agents"}"
// @Router /api/v1/projects/{project_id}/agents [get]
func GetAgents(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	var agents []models.Agent
	DB.Preload("Stage").Preload("Environment").
		Where("project_id = ? AND is_archived != ?", projectID, true).Find(&agents)

	// Convert to response
	var agentsResponse []agentResponse
	for _, agent := range agents {
		agentsResponse = append(agentsResponse, agentResponse{
			ID:            agent.ID,
			ProjectID:     agent.ProjectID,
			Name:          agent.Name,
			StageID:       agent.StageID,
			Stage:         agent.Stage,
			EnvironmentID: agent.EnvironmentID,
			Environment:   agent.Environment,
			WorkflowName:  agent.WorkflowName,
		})
	}
	c.JSON(http.StatusOK, gin.H{"agents": agentsResponse})
}

// GetAgentDetail godoc
// @Summary Get agent detail
// @Description Get agent detail
// @Tags Project Detail / Agents
// @Accept json
// @Produce json
// @Param agent_id path string true "Agent ID"
// @Param project_id path int true "Project ID"
// @Success 200 {object} models.Agent
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get agent by ID"}"
// @Router /api/v1/projects/{project_id}/agents/{agent_id} [get]
func GetAgentDetail(c *gin.Context) {
	agentID := c.Param("agent_id")
	projectID, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	var agent models.Agent
	result := DB.
		Preload("Stage").
		Preload("Environment").
		First(&agent, agentID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get agent by ID"})
		return
	}
	if agent.ProjectID != uint(projectID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent does not belong to project"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"agent": agent})
}

// GetProjectArchivedAgents godoc
// @Summary Get archived agents of project
// @Description Get archived agents of project
// @Tags Project Detail / Agents
// @Accept json
// @Produce json
// @Param project_id path int true "Project ID"
// @Success 200 {array}	 models.Agent
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get agents"}"
// @Router /api/v1/projects/{project_id}/agents/archived [get]
func GetArchivedAgents(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	var agents []models.Agent
	DB.Preload("Stage").Preload("Environment").
		Where("project_id = ? AND is_archived = ?", projectID, true).Find(&agents)

	// Convert to response
	var agentsResponse []agentResponse
	for _, agent := range agents {
		agentsResponse = append(agentsResponse, agentResponse{
			ID:            agent.ID,
			ProjectID:     agent.ProjectID,
			Name:          agent.Name,
			StageID:       agent.StageID,
			Stage:         agent.Stage,
			EnvironmentID: agent.EnvironmentID,
			Environment:   agent.Environment,
			WorkflowName:  agent.WorkflowName,
		})
	}
	c.JSON(http.StatusOK, gin.H{"agents": agentsResponse})
}

// ArchiveAgent godoc
// @Summary Archive agent
// @Description Archive agent
// @Tags Project Detail / Agents
// @Accept json
// @Produce json
// @Param agent_id path string true "Agent ID"
// @Param project_id path int true "Project ID"
// @Success 200 string {string} json "{"message": "Agent archived"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to archive agent"}"
// @Router /api/v1/projects/{project_id}/agents/{agent_id}/archive [put]
func ArchiveAgent(c *gin.Context) {
	agentID := c.Param("agent_id")
	projectID, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	var agent models.Agent
	result := DB.First(&agent, agentID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get agent by ID"})
		return
	}
	if agent.ProjectID != uint(projectID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent does not belong to project"})
		return
	}
	if agent.IsArchived {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent is already archived"})
		return
	}
	agent.IsArchived = true
	DB.Save(&agent)
	c.JSON(http.StatusOK, gin.H{"message": "Agent archived"})
}

// UnarchiveAgent godoc
// @Summary Unarchive agent
// @Description Unarchive agent
// @Tags Project Detail / Agents
// @Accept json
// @Produce json
// @Param agent_id path string true "Agent ID"
// @Param project_id path int true "Project ID"
// @Success 200 string {string} json "{"message": "Agent unarchived"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to unarchive agent"}"
// @Router /api/v1/projects/{project_id}/agents/{agent_id}/unarchive [put]
func RestoreAgent(c *gin.Context) {
	agentID := c.Param("agent_id")
	projectID, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	var agent models.Agent
	result := DB.First(&agent, agentID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get agent by ID"})
		return
	}
	if agent.ProjectID != uint(projectID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent does not belong to project"})
		return
	}
	if !agent.IsArchived {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent is not archived"})
		return
	}
	agent.IsArchived = false
	DB.Save(&agent)
	c.JSON(http.StatusOK, gin.H{"message": "Agent restored"})
}

type agentRequestBody struct {
	Name         string `json:"name" binding:"required"`
	Stage        string `json:"stage" binding:"required"`
	Environment  string `json:"environment" binding:"required"`
	WorkflowName string `json:"workflow_name" binding:"required"`
	Description  string `json:"description"`
}

// CreateNewAgent godoc
// @Summary Create new agent
// @Description Create new agent
// @Tags Project Detail / Agents
// @Accept json
// @Produce json
// @Param project_id path int true "Project ID"
// @Param Agent body models.Agent true "Agent"
// @Success 200 string {string} json "{"message": "Agent created"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to create agent"}"
// @Router /api/v1/projects/{project_id}/agents [post]
func CreateNewAgent(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	var agent agentRequestBody
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// find stage and environment by name agent provided
	var stage models.Stage
	result := DB.Where("name = ?", agent.Stage).First(&stage)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get stage by name"})
		return
	}
	var environment models.Environment
	result = DB.Where("name = ?", agent.Environment).First(&environment)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get environment by name"})
		return
	}
	// create new agent
	newAgent := models.Agent{
		ProjectID:     uint(projectID),
		Name:          agent.Name,
		StageID:       stage.ID,
		EnvironmentID: environment.ID,
		WorkflowName:  agent.WorkflowName,
		Description:   agent.Description,
		IsArchived:    false,
		ArchivedBy:    "",
		ArchivedAt:    time.Time{},
	}
	if err := DB.Create(&newAgent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create agent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent created"})
}

type requestAuthAgentBody struct {
	ApiToken string `json:"api_token" binding:"required"`
}

func agentLog(agent models.Agent, project models.Project, action string, message string, responseStatusCode int, latency time.Duration) {
	log := models.AgentLog{
		AgentID:        agent.ID,
		Agent:          agent,
		Action:         action,
		ProjectID:      project.ID,
		Project:        project,
		Message:        message,
		ResponseStatus: responseStatusCode,
		Latency:        int(latency.Milliseconds()),
	}
	DB.Create(&log)
}

// GetParameterByAuthAgent godoc
// @Summary Get parameter by auth agent
// @Description Get parameter by auth agent
// @Tags Agents
// @Accept json
// @Produce json
// @Param requestAuthAgentBody body controllers.requestAuthAgentBody true "Request Auth Agent Body"
// @Success 200 string {string} json "{"message": "Parameter retrieved"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to retrieve parameter"}"
// @Router /api/v1/agents/auth-parameters [post]
func GetParameterByAuthAgent(c *gin.Context) {
	var reqBody requestAuthAgentBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var agent models.Agent
	result := DB.Where("api_token = ?", reqBody.ApiToken).First(&agent)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get agent by API token"})
		return
	}
	startTime := time.Now()
	var project models.Project
	if err := DB.
		Preload("LatestVersion").
		Preload("LatestVersion.Parameters", "stage_id = ? AND environment_id = ? AND is_archived = ? ", agent.StageID, agent.EnvironmentID, false).
		First(&project, agent.ProjectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get project by agent"})
		return
	}
	latency := time.Since(startTime)
	agentLog(agent, project, "Get Params", "Succeed: Applied previous updated param", http.StatusOK, latency)
	c.JSON(http.StatusOK, gin.H{"parameters": project.LatestVersion.Parameters})
}

// RerunWorkFlowByAgent godoc
// @Summary Rerun workflow by agent
// @Description Rerun workflow by agent
// @Tags Agents
// @Accept json
// @Produce json
// @Param agent_id path string true "Agent ID"
// @Success 200 string {string} json "{"message": "rerun workflow by agent"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to rerun workflow by agent"}"
// @Router /api/v1/agents/{agent_id}/rerun-workflow [post]
func RerunWorkFlowByAgent(c *gin.Context) {
	agent_id := c.Param("agent_id")
	var agent models.Agent
	result := DB.
		Preload("Stage").
		Preload("Environment").
		First(&agent, agent_id)
	if result.Error != nil {
		log.Println("Failed to get agent by ID")
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get agent by ID"})
		return
	}

	var project models.Project
	result = DB.First(&project, agent.ProjectID)
	if result.Error != nil {
		log.Println("Failed to get project by agent ID")
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get project by agent ID"})
		return
	}
	githubRepository, err := github.ParseRepoURL(project.RepoURL)
	if err != nil {
		log.Println("Failed to parse repo URL")
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to parse repo URL"})
		return
	}
	startTime := time.Now()
	responseStatusCode, responseBodyMessage, err := github.RerunWorkFlow(githubRepository.Owner, githubRepository.Name, agent.WorkflowName, project.RepoApiToken)
	latency := time.Since(startTime)

	if responseStatusCode == http.StatusForbidden { // if workflow is already running
		responseStatusCode = http.StatusAccepted
		responseBodyMessage = "CICD is already running"
	}
	if responseStatusCode == http.StatusCreated { // if workflow is rerun
		responseBodyMessage = "CICD is starting rerun"
	}
	if err != nil {
		rerunLog(project.ID, agent.ID, responseStatusCode, responseBodyMessage, responseStatusCode, latency)
		c.JSON(responseStatusCode, gin.H{
			"latency": latency.String(),
			"status":  responseStatusCode,
			"message": responseBodyMessage,
		})
		return
	}
	rerunLog(project.ID, agent.ID, responseStatusCode, responseBodyMessage, responseStatusCode, latency)
	c.JSON(responseStatusCode, gin.H{
		"latency": latency.String(),
		"status":  responseStatusCode,
		"message": responseBodyMessage,
	})
}

func rerunLog(projectID uint, agentID uint, responseStatus int, message string, cicdResponseCode int, latency time.Duration) {
	log := models.AgentLog{
		ProjectID:      projectID,
		AgentID:        agentID,
		ResponseStatus: responseStatus,
		Action:         "Rerun Workflow",
		Latency:        int(latency.Milliseconds()),
		// Message:        message,
	}
	switch cicdResponseCode {
	case 201:
		log.Message = "Created: CICD is starting rerun"
	case 202:
		log.Message = "Accepted: CICD is already running"
	case 401:
		log.Message = fmt.Sprintf("Unauthorized: %s", message)
	case 404:
		log.Message = fmt.Sprintf("Not Found: %s", message)
	case 500:
		log.Message = "Internal Server Error"
	}

	DB.Create(&log)
}
