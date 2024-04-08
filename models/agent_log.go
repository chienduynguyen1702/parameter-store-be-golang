package models

import (
	"gorm.io/gorm"
)

type AgentLog struct {
	gorm.Model
	AgentID uint   `gorm:"foreignKey:AgentID" json:"agent_id"`
	Agent   Agent  `json:"agent"`
	Path    string `gorm:"type:varchar(100);not null" json:"path"`
	Status  int    `gorm:"not null" json:"status"`
	Latency int    `gorm:"not null" json:"latency"`
}
