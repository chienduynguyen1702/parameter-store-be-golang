package models

import "gorm.io/gorm"

type Workflow struct {
	gorm.Model
	IsActivated      bool          `gorm:"default:false" json:"is_activated"`
	ProjectID        uint          `gorm:"foreignKey:ProjectID" json:"project_id"`
	IsUpdatedLastest bool          `gorm:"default:false" json:"is_updated_lastest"`
	Logs             []WorkflowLog `gorm:"one2many:workflow_logs;" json:"logs"`
}
