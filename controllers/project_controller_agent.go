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

type agentResponse struct {
	ID            uint `json:"id"`
	ProjectID     uint
	Name          string `gorm:"type:varchar(100);not null" json:"name"`
	StageID       uint   `gorm:"foreignKey:StageID;not null" json:"stage_id"`
	Stage         models.Stage
	EnvironmentID uint `gorm:"foreignKey:EnvironmentID;not null" json:"environment_id"`
	Environment   models.Environment
	WorkflowName  string `gorm:"type:varchar(100);not null" json:"workflow_name"`
	Description   string `gorm:"type:varchar(100);not null" json:"description"`
	LastUsedAt    time.Time
	ArchivedAt    time.Time `json:"archived_at"`
	ArchivedBy    string    `json:"archived_by"`
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
	DB.Preload("Stage").Preload("Environment").Preload("Workflow").
		Where("project_id = ? AND is_archived != ?", projectID, true).Find(&agents)
	totalAgents := len(agents)
	page := c.Query("page")
	limit := c.Query("limit")

	var agentsResponse []agentResponse
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
		paginatedAgents := paginationDataAgent(agents, pageInt, limitInt)
		for _, agent := range paginatedAgents {
			agentsResponse = append(agentsResponse, agentResponse{
				ID:            agent.ID,
				ProjectID:     agent.ProjectID,
				Name:          agent.Name,
				StageID:       agent.StageID,
				Stage:         agent.Stage,
				Description:   agent.Description,
				EnvironmentID: agent.EnvironmentID,
				Environment:   agent.Environment,
				WorkflowName:  agent.Workflow.Name,
				LastUsedAt:    agent.LastUsedAt,
			})
		}
	}
	// Convert to response
	c.JSON(http.StatusOK, gin.H{
		"agents": agentsResponse,
		"total":  totalAgents,
	})
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
			ArchivedAt:    agent.ArchivedAt,
			ArchivedBy:    agent.ArchivedBy,
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

	// get username from context
	user, exist := c.Get("user")
	if !exist {
		log.Println("Failed to get user from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}

	agent.IsArchived = true
	agent.ArchivedAt = time.Now()
	agent.ArchivedBy = user.(models.User).Username
	if err := DB.Save(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to archive agent"})
		return
	}
	projectLogByUser(uint(projectID), "Archive Agent", "Succeed: Agent archived", http.StatusOK, time.Since(time.Now()), 0)
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
	projectLogByUser(uint(projectID), "Restore Agent", "Succeed: Agent restored", http.StatusOK, time.Since(time.Now()), 0)
	c.JSON(http.StatusOK, gin.H{"message": "Agent restored"})
}

type agentRequestBody struct {
	Name          string `json:"name" binding:"required"`
	Stage         string `json:"stage"  binding:"required"`
	StageID       uint   `json:"stage_id"`
	Environment   string `json:"environment"  binding:"required"`
	EnvironmentID uint   `json:"environment_id"`
	WorkflowName  string `json:"workflow_name" binding:"required"`
	Description   string `json:"description" binding:"required"`
}

// CreateNewAgent godoc
// @Summary Create new agent
// @Description Create new agent
// @Tags Project Detail / Agents
// @Accept json
// @Produce json
// @Param project_id path int true "Project ID"
// @Param Agent body controllers.agentRequestBody true "Agent"
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
	//debug
	// log.Println(agent)
	// find stage and environment in project by projectID
	project := models.Project{}
	DB.Preload("Stages").Preload("Environments").Preload("Workflows").First(&project, projectID)
	// validate workflow name
	if err := github.ValidateWorkflowName(agent.WorkflowName, project.RepoURL, project.RepoApiToken); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	//find stage id and environment id workfow id
	for _, stage := range project.Stages {
		if stage.Name == agent.Stage {
			agent.StageID = stage.ID
			break
		}
	}
	for _, environment := range project.Environments {
		if environment.Name == agent.Environment {
			agent.EnvironmentID = environment.ID
			break
		}
	}
	var findingWorkflowID uint
	for _, workflow := range project.Workflows {
		if workflow.Name == agent.WorkflowName {
			findingWorkflowID = uint(workflow.WorkflowID)
			break
		}
	}
	// create new agent
	newAgent := models.Agent{
		ProjectID:     uint(projectID),
		Name:          agent.Name,
		StageID:       agent.StageID,
		EnvironmentID: agent.EnvironmentID,
		WorkflowID:    findingWorkflowID,
		WorkflowName:  agent.WorkflowName,
		Description:   agent.Description,
		IsArchived:    false,
		ArchivedBy:    "",
		ArchivedAt:    time.Time{},
		// APIToken:      ApiTokenForAgent,
	}
	if err := DB.Create(&newAgent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create agent"})
		return
	}
	ApiTokenForAgent := GenerateTokenForAgent(strconv.Itoa(int(newAgent.ID)), strconv.Itoa(int(project.OrganizationID)))
	newAgent.APIToken = ApiTokenForAgent
	DB.Save(&newAgent)
	projectLogByUser(uint(projectID), "Create Agent", "Succeed: Agent created", http.StatusCreated, time.Since(time.Now()), 0)
	c.JSON(http.StatusOK, gin.H{
		"message":   "Agent created",
		"api_token": newAgent.APIToken,
	})
}

// UpdateAgent godoc
// @Summary Update agent
// @Description Update agent
// @Tags Project Detail / Agents
// @Accept json
// @Produce json
// @Param agent_id path string true "Agent ID"
// @Param project_id path int true "Project ID"
// @Param Agent body controllers.agentRequestBody true "Agent"
// @Success 200 string {string} json "{"message": "Agent updated"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to update agent"}"
// @Router /api/v1/projects/{project_id}/agents/{agent_id} [put]
func UpdateAgent(c *gin.Context) {
	agentID := c.Param("agent_id")
	projectID, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	var agent agentRequestBody
	if err := c.ShouldBindJSON(&agent); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// find stage and environment in project by projectID
	project := models.Project{}
	DB.Preload("Stages").Preload("Environments").First(&project, projectID)
	// validate workflow name
	if err := github.ValidateWorkflowName(agent.WorkflowName, project.RepoURL, project.RepoApiToken); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	//find stage id and environment id
	for _, stage := range project.Stages {
		if stage.Name == agent.Stage {
			agent.StageID = stage.ID
			break
		}
	}
	for _, environment := range project.Environments {
		if environment.Name == agent.Environment {
			agent.EnvironmentID = environment.ID
			break
		}
	}
	var findingWorkflowID uint
	for _, workflow := range project.Workflows {
		if workflow.Name == agent.WorkflowName {
			findingWorkflowID = uint(workflow.WorkflowID)
			break
		}
	}

	// update agent
	agentUpdate := models.Agent{
		Name:          agent.Name,
		StageID:       agent.StageID,
		EnvironmentID: agent.EnvironmentID,
		WorkflowID:    findingWorkflowID,
		WorkflowName:  agent.WorkflowName,
		Description:   agent.Description,
	}
	if err := DB.Model(&models.Agent{}).Where("id = ?", agentID).Updates(agentUpdate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update agent"})
		return
	}
	projectLogByUser(uint(projectID), "Update Agent", "Succeed: Agent updated", http.StatusOK, time.Since(time.Now()), 0)
	c.JSON(http.StatusOK, gin.H{"message": "Agent updated"})
}

type requestAuthAgentBody struct {
	ApiToken string `json:"api_token" binding:"required"`
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Failed to authorized agent by API token, please check $PARAMETER_STORE_TOKEN",
		})
		return
	}
	agent.LastUsedAt = time.Now()
	DB.Save(&agent)
	startTime := time.Now()
	var project models.Project
	if err := DB.
		Preload("LatestVersion").
		Preload("LatestVersion.Parameters", "stage_id = ? AND environment_id = ? AND is_archived = ? ", agent.StageID, agent.EnvironmentID, false).
		First(&project, agent.ProjectID).Error; err != nil {

		agentLog(agent, project, "Get Parameter", "Failed to get project by agent", http.StatusNotFound, time.Since(startTime))
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Failed to get project by agent",
		})

		return
	}
	for _, parameter := range project.LatestVersion.Parameters {
		parameter.IsApplied = true
		DB.Save(&parameter)
	}
	latency := time.Since(startTime)
	agentLog(agent, project, "Get Parameter", "Succeed: Parameter retrieved", http.StatusOK, latency)
	c.JSON(http.StatusOK, gin.H{
		"status":     http.StatusOK,
		"message":    "Parameter retrieved",
		"parameters": project.LatestVersion.Parameters,
	})
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
