package controllers

import (
	"log"
	"net/http"
	"parameter-store-be/models"
	"parameter-store-be/modules/github"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
			// log.Println("workflow ID  : ", workflow.ID)
			// log.Println("workflow Name: ", workflow.Name)
			// log.Println("Lastest workflow run ID: ", lastestWorkflowRunID)
			// log.Println("Last attempt number  : ", lastAttemptNumber)
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
				// update the workflow
				wf.Name = workflow.Name
				wf.AttemptNumber = lastAttemptNumber
				wf.LastWorkflowRunID = lastestWorkflowRunID
				wf.State = workflow.State
				wf.Path = workflow.Path

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

// GetWorkflowLogs is a function to get workflow history
// @Summary Get workflow history
// @Description Get workflow history
// @Tags Project Detail / Workflows
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param workflow_id path string true "Workflow ID"
// @Success 200 string {string} json "{"workflow": "workflow"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get workflow history"}"
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/workflows/{workflow_id}/logs [get]
func GetWorkflowLogs(c *gin.Context) {
	projectID := c.Param("project_id")
	workflowID := c.Param("workflow_id")
	// fmt.Println("Project ID: ", projectID)
	// fmt.Println("Workflow ID: ", workflowID)
	var prj models.Project
	if err := DB.
		Preload("Workflows", "workflow_id = ? ", workflowID).
		Preload("Workflows.Logs",
			func(db *gorm.DB) *gorm.DB { // order by workflow_log created_at desc
				db = db.Order("workflow_logs.created_at desc")
				return db
			}).First(&prj, projectID).Error; err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get project"})
		return
	}

	// debug
	// fmt.Println("prj: ", prj)

	var workflow models.Workflow
	var workflowLogs []models.WorkflowLog
	if len(prj.Workflows) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	} else { // assign workflow and workflowLogs
		workflow = prj.Workflows[0]
		workflowLogs = prj.Workflows[0].Logs
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"workflows": workflow,
			"logs":      workflowLogs,
		},
	})
}

// GetDiffParameterInWorkflowLog is a function to get diff parameter in workflow log, pulled by agent
// @Summary Get diff parameter in workflow log
// @Description Get diff parameter pulled by agent in 2 nearest workflow log
// @Tags Project Detail / Workflows
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param workflow_id path string true "Workflow ID"
// @Param workflow_log_id path string true "Workflow Log ID"
// @Success 200 string {string} json "{"workflow": "workflow"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get diff parameter in workflow log"}"
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/workflows/{workflow_id}/logs/{workflow_log_id}/diff-parameter [get]
func GetDiffParameterInWorkflowLog(c *gin.Context) {
	projectID := c.Param("project_id")
	workflowID := c.Param("workflow_id")
	workflowLogID := c.Param("workflow_log_id")
	var project models.Project
	if err := DB.
		Preload("Workflows", "workflow_id = ? ", workflowID).
		Preload("Workflows.Logs", "id = ?", workflowLogID).
		Preload("Workflows.Logs.AgentLogs").
		Preload("Workflows.Logs.AgentLogs.Agent").
		Preload("Workflows.Logs.AgentLogs.Agent.Stage").
		Preload("Workflows.Logs.AgentLogs.AgentPullParameterLog").
		First(&project, projectID).Error; err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get project"})
	}
	var curentWorkflowLog models.WorkflowLog
	if len(project.Workflows) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	} else if len(project.Workflows[0].Logs) == 0 { // assign workflow and workflowLogs
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow log not found"})
		return
	} else {
		curentWorkflowLog = project.Workflows[0].Logs[0]
	}

	var previousWorkflowLog models.WorkflowLog

	err := DB.Where("workflow_id = ? AND ((workflow_run_id < ?) OR (workflow_run_id = ? AND attempt_number < ?))",
		workflowID, curentWorkflowLog.WorkflowRunId, curentWorkflowLog.WorkflowRunId, curentWorkflowLog.AttemptNumber).
		Order("workflow_run_id DESC, attempt_number DESC").
		Preload("AgentLogs").
		Preload("AgentLogs.Agent").
		Preload("AgentLogs.Agent.Stage").
		Preload("AgentLogs.AgentPullParameterLog").
		First(&previousWorkflowLog).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Previous workflow log not found"})
		return
	}
	diff := getDiffParameterBetween2WorkflowLogs(curentWorkflowLog, previousWorkflowLog)

	c.JSON(http.StatusOK, gin.H{
		// "current":  curentWorkflowLog,
		// "previous": previousWorkflowLog,
		"diff": diff,
	})

}

type ParameterDiffBetweenWorkflow struct {
	Stages []stagesBetweenWorkflow `json:"stages"`
}
type stagesBetweenWorkflow struct {
	ID         uint                             `json:"id"`
	Name       string                           `json:"name"`
	Parameters []parameterChangeBetweenWorkflow `json:"parameters"`
}

type parameterChangeBetweenWorkflow struct {
	Name          string `json:"name"`
	CurrentValue  string `json:"current_value"`
	PreviousValue string `json:"previous_value"`
}

func getDiffParameterBetween2WorkflowLogs(first, second models.WorkflowLog) ParameterDiffBetweenWorkflow {
	var stages []stagesBetweenWorkflow
	stageMap := make(map[string]*stagesBetweenWorkflow)

	// Process the first workflow log
	for _, agentLog := range first.AgentLogs {
		stageName := agentLog.Agent.Stage.Name
		stageID := agentLog.Agent.Stage.ID
		if _, exists := stageMap[stageName]; !exists {
			stageMap[stageName] = &stagesBetweenWorkflow{
				ID:         stageID,
				Name:       stageName,
				Parameters: []parameterChangeBetweenWorkflow{},
			}
		}
		for _, parameter := range agentLog.AgentPullParameterLog {
			param := parameterChangeBetweenWorkflow{
				Name:          parameter.ParameterName,
				CurrentValue:  parameter.ParameterValue,
				PreviousValue: "",
			}
			stageMap[stageName].Parameters = append(stageMap[stageName].Parameters, param)
		}
	}

	// Process the second workflow log
	for _, agentLog := range second.AgentLogs {
		stageName := agentLog.Agent.Stage.Name
		if _, exists := stageMap[stageName]; !exists {
			stageMap[stageName] = &stagesBetweenWorkflow{
				Name:       stageName,
				Parameters: []parameterChangeBetweenWorkflow{},
			}
		}
		for _, parameter := range agentLog.AgentPullParameterLog {
			updated := false
			for i, param := range stageMap[stageName].Parameters {
				if param.Name == parameter.ParameterName {
					stageMap[stageName].Parameters[i].PreviousValue = parameter.ParameterValue
					updated = true
					break
				}
			}
			if !updated {
				param := parameterChangeBetweenWorkflow{
					Name:          parameter.ParameterName,
					CurrentValue:  "",
					PreviousValue: parameter.ParameterValue,
				}
				stageMap[stageName].Parameters = append(stageMap[stageName].Parameters, param)
			}
		}
	}

	// Convert the map to a slice
	for _, stage := range stageMap {
		stages = append(stages, *stage)
	}
	// Sort the stages slice by stage name
	sort.Slice(stages, func(i, j int) bool {
		return stages[i].ID < stages[j].ID
	})

	parameterDiff := ParameterDiffBetweenWorkflow{
		Stages: stages,
	}
	return parameterDiff
}
