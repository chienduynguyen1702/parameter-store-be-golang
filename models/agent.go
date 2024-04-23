package models

import (
	"time"

	"gorm.io/gorm"
)

type Agent struct {
	gorm.Model
	ProjectID     uint
	Name          string    `gorm:"type:varchar(100);not null" json:"name"`
	APIToken      string    `gorm:"type:varchar(100)" json:"api_token"`
	StageID       uint      `gorm:"foreignKey:StageID;not null" json:"stage_id"`
	EnvironmentID uint      `gorm:"foreignKey:EnvironmentID;not null" json:"environment_id"`
	LastUsedAt    time.Time `gorm:"type:timestamp;" json:"last_used_at"`
	WorkflowName  string    `gorm:"type:varchar(100);not null" json:"workflow_name"`
	WorkflowID    uint      `gorm:"foreignKey:WorkflowID" json:"workflow_id"`
	Description   string    `gorm:"type:text" json:"description"`
	IsArchived    bool      `gorm:"default:false" json:"is_archived"`
	ArchivedBy    string    `gorm:"foreignKey:ArchivedBy" json:"archived_by"` // foreign key to user model
	ArchivedAt    time.Time `gorm:"type:timestamp;" json:"archived_at"`

	// Workflow    Workflow
	Stage       Stage
	Environment Environment
}
