package models

import (
	"time"

	"gorm.io/gorm"
)

type WorkflowLog struct {
	gorm.Model
	WorkflowID    uint      `json:"workflow_id"`
	WorkflowRunId uint      `json:"workflow_run_id"`
	AttemptNumber int       `json:"attempt_number"`
	State         string    `json:"state"`
	StartedAt     time.Time `json:"started_at"`
	Duration      int       `json:"duration"`
	// ProjectLogID  uint       `json:"project_log_id"`
	// ProjectLog    ProjectLog `gorm:"foreignKey:ProjectLogID" json:"project_log"`
}
