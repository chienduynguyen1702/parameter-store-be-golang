package models

import (
	"gorm.io/gorm"
)

type ProjectLog struct {
	gorm.Model
	// AgentID        uint   `gorm:"foreignKey:AgentID" json:"agent_id"`
	// IsByUser       bool   `gorm:"default:false" json:"is_by_user"`
	UserID         uint   `gorm:"foreignKey:UserID" json:"user_id"`
	Action         string `gorm:"type:varchar(100);not null" json:"action"`
	ProjectID      uint   `gorm:"foreignKey:ProjectID" json:"project_id"`
	Path           string `gorm:"type:varchar(100);not null" json:"path"`
	ResponseStatus int    `gorm:"not null" json:"response_status"`
	Message        string `gorm:"type:text" json:"message"`
	Latency        int    `gorm:"not null" json:"latency"`

	// Agent   Agent   `json:"agent"`
	User    User    `json:"user"`
	Project Project `json:"project"`
}
