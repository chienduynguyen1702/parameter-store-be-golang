package controllers

import (
	"log"
	"net/http"
	"parameter-store-be/models"
	"parameter-store-be/modules/github"
	"time"

	"github.com/gin-gonic/gin"
)

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
