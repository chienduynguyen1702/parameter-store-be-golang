package models

import (
	"time"

	"gorm.io/gorm"
)

type WorkflowLog struct {
	gorm.Model
	StartAt    time.Time `json:"start_at"`
	Latency    time.Time `json:"latency"`
	WorkflowID uint      `gorm:"foreignKey:WorkflowID" json:"workflow_id"`
}
