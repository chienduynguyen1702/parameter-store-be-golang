package controllers

import (
	"log"
	"net/http"
	"parameter-store-be/models"
	"parameter-store-be/modules/github"

	"github.com/gin-gonic/gin"
)

// GetProjectWorkflows is a function to get project workflows
// @Summary Get project workflows
// @Description Get project workflows
// @Tags Project Detail / Workflows
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"workflow": "workflow"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get project workflows"}"
// @Router /api/v1/projects/{project_id}/workflows [get]
func GetProjectWorkflows(c *gin.Context) {
	projectID := c.Param("project_id")
	var project models.Project
	result := DB.First(&project, projectID)
	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}
	listWorkflows, err := github.GetWorkflows(project.RepoURL, project.RepoApiToken)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get workflows"})
		return
	}
	// debug
	// fmt.Println("List workflows in GetProjectWorkflows : ", listWorkflows, "\n")
	// Save the listWorkflows of project back to the database
	for _, workflow := range listWorkflows.Workflows {
		repo, err := github.ParseRepoURL(project.RepoURL)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse repository URL"})
			return
		}

		_, _, lastAttemptNumber, err := github.GetLastAttemptNumberOfWorkflowRun(repo.Owner, repo.Name, project.RepoApiToken, workflow.Name)

		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get last attempt number"})
			return
		}
		// log.Println("Last attempt number: ", lastAttemptNumber)
		var wf models.Workflow
		result := DB.Where("workflow_id = ? AND project_id = ?", workflow.ID, project.ID).First(&wf)
		if result.RowsAffected == 0 {
			// log.Println("Workflow id: ", workflow.ID)
			wf = models.Workflow{
				WorkflowID:    uint(workflow.ID),
				Name:          workflow.Name,
				Path:          workflow.Path,
				ProjectID:     project.ID,
				State:         workflow.State,
				AttemptNumber: lastAttemptNumber,
			}
			DB.Create(&wf)
		} else {
			wf.AttemptNumber = lastAttemptNumber
			DB.Save(&wf)
		}
	}
	// preload workflows from the database using the project ID
	DB.Model(&project).Preload("Workflows").Find(&project)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"repo_url":  project.RepoURL,
			"workflows": project.Workflows,
		},
	})
}
