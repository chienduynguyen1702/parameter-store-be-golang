package controllers

import (
	"net/http"
	"parameter-store-be/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetProjectAgents godoc
// @Summary Get agents of project
// @Description Get agents of project
// @Tags Project Detail / Agent
// @Accept json
// @Produce json
// @Param project_id path int true "Project ID"
// @Success 200 {array}	 models.Agent
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get agents"}"
// @Router /api/v1/project/{project_id}/agent [get]
func GetProjectAgents(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	var agents []models.Agent
	DB.Where("project_id = ?", projectID).Find(&agents)
	c.JSON(http.StatusOK, gin.H{"agents": agents})
}

// CreateNewAgent godoc
// @Summary Create new agent
// @Description Create new agent
// @Tags Project Detail / Agent
// @Accept json
// @Produce json
// @Param project_id path int true "Project ID"
// @Param Agent body models.Agent true "Agent"
// @Success 200 string {string} json "{"message": "Agent created"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to create agent"}"
// @Router /api/v1/project/{project_id}/agent [post]
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
