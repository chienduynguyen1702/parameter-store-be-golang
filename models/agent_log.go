package models

import "time"

type AgentLog struct {
	ID      uint `gorm:"primaryKey"`
	AgentID uint
	Agent   Agent
	Path    string `gorm:"type:varchar(100);not null"`
	Time    *time.Time
	Status  int `gorm:"not null"`
	Latency int
}
