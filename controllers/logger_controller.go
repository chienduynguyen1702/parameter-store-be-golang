package controllers

import (
	"fmt"
	"parameter-store-be/models"
	"time"
)

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

func projectLogByUser(projectID uint, action string, message string, responseStatusCode int, latency time.Duration, userID uint) {
	log := models.ProjectLog{
		// IsByUser:       true,
		// AgentID:        0,
		UserID:         userID,
		Action:         action,
		ProjectID:      projectID,
		Message:        message,
		ResponseStatus: responseStatusCode,
		Latency:        int(latency.Milliseconds()),
	}
	DB.Create(&log)
}

// func projectLogByAgent(projectID uint, agentID uint, action string, message string, responseStatusCode int, latency time.Duration) {
// 	log := models.ProjectLog{
// 		IsByUser:       false,
// 		UserID:         0,
// 		// AgentID:        agentID,
// 		Action:         action,
// 		ProjectID:      projectID,
// 		Message:        message,
// 		ResponseStatus: responseStatusCode,
// 		Latency:        int(latency.Milliseconds()),
// 	}
// 	DB.Create(&log)

// }

func workflowLog(workflowID uint, workflowRunId uint, attemptNumber int) {
	// find workflow by ID
	var workflow models.Workflow
	DB.Where("workflow_id = ?", workflowID).First(&workflow)
	workflow.IsUpdatedLastest = false
	workflow.AttemptNumber = attemptNumber
	DB.Save(&workflow)

	// create workflow log
	log := models.WorkflowLog{
		WorkflowID:    workflowID,
		WorkflowRunId: workflowRunId,
		AttemptNumber: attemptNumber,
		ProjectID:     workflow.ProjectID,
	}
	DB.Create(&log)
}
