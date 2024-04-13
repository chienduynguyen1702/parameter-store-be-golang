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
}

// GetProjectAgents godoc
// @Summary Get agents of project
// @Description Get agents of project
// @Tags Project Detail / Agents
// @Accept json
// @Produce json
// @Param project_id path int true "Project ID"
// @Success 200 {array}	 models.Agent
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get agents"}"
// @Router /api/v1/projects/{project_id}/agents [get]
func GetProjectAgents(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	var agents []models.Agent
	DB.Preload("Stage").Preload("Environment").
		Where("project_id = ?", projectID).Find(&agents)

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
	var agent models.Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	agent.ProjectID = uint(projectID)
	DB.Create(&agent)
	c.JSON(http.StatusOK, gin.H{"message": "Agent created"})
}

// RerunWorkFlowByAgent godoc
// @Summary Rerun workflow by agent
// @Description Rerun workflow by agent
// @Tags Project Detail / Agents
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
	responseStatusCode, err := github.RerunWorkFlow(githubRepository.Owner, githubRepository.Name, agent.WorkflowName, project.RepoApiToken)
	latency := time.Since(startTime)
	if err != nil {
		// fmt.Println(err)
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"latency": latency.String(),
			"status":  responseStatusCode,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"latency": latency.String(),
		"status":  responseStatusCode,
		"message": "rerun workflow by agent successfully",
	})
}
