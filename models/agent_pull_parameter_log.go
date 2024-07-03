package models

import "gorm.io/gorm"

type AgentPullParameterLog struct {
	gorm.Model
	AgentLogID     uint   `gorm:"foreignKey:AgentLogID" json:"agent_log_id"`
	AgentID        uint   `gorm:"foreignKey:AgentID" json:"agent_id"`
	ParameterID    uint   `gorm:"foreignKey:ParameterID" json:"parameter_id"`
	ProjectID      uint   `gorm:"foreignKey:ProjectID" json:"project_id"`
	ParameterValue string `gorm:"type:text" json:"parameter_value"`
	ParameterName  string `gorm:"type:varchar(100);not null" json:"parameter_name"`
	StageID        uint   `gorm:"foreignKey:StageID" json:"stage_id"`
	EnvironmentID  uint   `gorm:"foreignKey:EnvironmentID" json:"environment_id"`

	Stage       Stage       `json:"stage"`
	Environment Environment `json:"environment"`
	Project     Project     `json:"project"`
	Parameter   Parameter   `json:"parameter"`
}
