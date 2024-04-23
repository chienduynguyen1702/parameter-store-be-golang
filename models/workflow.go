package models

import "gorm.io/gorm"

type Workflow struct {
	gorm.Model
	// ID            uint   `gorm:"primaryKey" json:"id"`
	WorkflowID    uint    `gorm:"primaryKey" json:"workflow_id"`
	Name          string `json:"name"`
	Path          string `json:"path"`
	State         string `json:"state"`
	AttemptNumber int    `json:"attempt_number"`
	// IsActivated      bool          `gorm:"default:false" json:"is_activated"`
	ProjectID        uint          `gorm:"foreignKey:ProjectID" json:"project_id"`
	IsUpdatedLastest bool          `gorm:"default:false" json:"is_updated_lastest"`
	Logs             []WorkflowLog ` json:"logs"`
}
