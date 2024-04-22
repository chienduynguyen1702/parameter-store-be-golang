package models

import "gorm.io/gorm"

type Workflow struct {
	gorm.Model
	WorkflowID int    `gorm:"primaryKey" json:"workflow_id"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	State      string `json:"state"`
	// IsActivated      bool          `gorm:"default:false" json:"is_activated"`
	ProjectID        uint          `gorm:"foreignKey:ProjectID" json:"project_id"`
	IsUpdatedLastest bool          `gorm:"default:false" json:"is_updated_lastest"`
	Logs             []WorkflowLog `gorm:"one2many:workflow_logs;" json:"logs"`
}
