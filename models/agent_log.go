package models

import (
	"gorm.io/gorm"
)

type AgentLog struct {
	gorm.Model
	AgentID uint
	Agent   Agent
	Path    string `gorm:"type:varchar(100);not null"`
	Status  int    `gorm:"not null"`
	Latency int
}
