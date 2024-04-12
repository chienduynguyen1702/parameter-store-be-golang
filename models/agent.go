package models

import (
	"gorm.io/gorm"
)

type Agent struct {
	gorm.Model
	ProjectID     uint
	Name          string `gorm:"type:varchar(100);not null" json:"name"`
	APIToken      string `gorm:"type:varchar(100)" json:"api_token"`
	StageID       uint   `gorm:"foreignKey:StageID;not null" json:"stage_id"`
	Stage         Stage
	EnvironmentID uint `gorm:"foreignKey:EnvironmentID;not null" json:"environment_id"`
	Environment   Environment
	WorkflowName  string `gorm:"type:varchar(100);not null" json:"workflow_name"`
}
