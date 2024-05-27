package controllers

import (
	"log"
	"net/http"
	"parameter-store-be/models"
	"parameter-store-be/modules/github"
	"strconv"

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
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/workflows [get]
func GetProjectWorkflows(c *gin.Context) {
	projectID := c.Param("project_id")
	var project models.Project
	result := DB.Preload("Workflows").First(&project, projectID)
	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get project"})
		return
	}
	go func() {
		listWorkflows, err := github.GetWorkflows(project.RepoURL, project.RepoApiToken)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get workflows"})
			return
		}
		// debug
		// fmt.Println("List workflows in GetProjectWorkflows : ", listWorkflows, "\n")
		// Save the listWorkflows of project back to the database
		for _, workflow := range listWorkflows.Workflows {
			repo, err := github.ParseRepoURL(project.RepoURL)
			if err != nil {
				log.Println(err.Error())
				c.JSON(http.StatusNotFound, gin.H{"error": "Failed to parse repository URL"})
				return
			}

			lastestWorkflowRunIDString, _, lastAttemptNumber, err := github.GetLastAttemptNumberOfWorkflowRun(repo.Owner, repo.Name, project.RepoApiToken, workflow.Name)

			if err != nil {
				log.Println(err.Error())
				c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get last attempt number"})
				return
			}
			// parse lastestWorkflowRunID to int
			lastestWorkflowRunID, err := strconv.Atoi(lastestWorkflowRunIDString)
			if err != nil {
				log.Println(err.Error())
				c.JSON(http.StatusNotFound, gin.H{"error": "Failed to parse lastestWorkflowRunID to int"})
			}

			// log.Println("Last attempt number: ", lastAttemptNumber)
			var wf models.Workflow
			result := DB.Where("workflow_id = ? AND project_id = ?", workflow.ID, project.ID).First(&wf)
			if result.RowsAffected == 0 {
				// log.Println("Workflow id: ", workflow.ID)
				wf = models.Workflow{
					WorkflowID:        uint(workflow.ID),
					Name:              workflow.Name,
					Path:              workflow.Path,
					ProjectID:         project.ID,
					State:             workflow.State,
					AttemptNumber:     lastAttemptNumber,
					LastWorkflowRunID: lastestWorkflowRunID,
				}
				DB.Create(&wf)
			} else {
				wf.AttemptNumber = lastAttemptNumber
				wf.LastWorkflowRunID = lastestWorkflowRunID
				DB.Save(&wf)
			}
		}
	}()
	// preload workflows from the database using the project ID
	DB.Model(&project).Preload("Workflows").Find(&project)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"repo_url":  project.RepoURL,
			"workflows": project.Workflows,
		},
	})
}

// GetWorkflowProcess is a function to get workflow process
// @Summary Get workflow process
// @Description Get workflow process
// @Tags Project Detail / Workflows
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param workflow_id path string true "Workflow ID"
// @Success 200 string {string} json "{"workflow": "workflow"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get workflow process"}"
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/workflows/{workflow_id}/run [get]
func GetWorkflowProcess(c *gin.Context) {
	projectID := c.Param("project_id")
	workflowID := c.Param("workflow_id")
	// fmt.Println("Project ID: ", projectID)
	// fmt.Println("Workflow ID: ", workflowID)
	var prj models.Project
	prjResult := DB.First(&prj, projectID)
	if prjResult.Error != nil {
		log.Println(prjResult.Error)
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get project"})
		return
	}
	prjRepo, err := github.ParseRepoURL(prj.RepoURL)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to parse repository URL"})
		return
	}

	var workflow models.Workflow
	workflowResult := DB.Where("workflow_id = ? AND project_id = ?", workflowID, projectID).First(&workflow)
	if workflowResult.Error != nil {
		log.Println(workflowResult.Error)
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get project"})
		return
	}
	workflowJobs, err := github.ListJobsForAWorkflowRun(prjRepo.Owner, prjRepo.Name, prj.RepoApiToken, workflow.LastWorkflowRunID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get workflow jobs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"workflows": workflow,
			"jobs":      workflowJobs,
		},
	})
}
