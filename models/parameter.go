package models

import (
	"gorm.io/gorm"
)

type Parameter struct {
	gorm.Model
	StageID       uint        `gorm:"foreignKey:StageID" json:"stage_id"`
	Stage         Stage       `json:"stage"`
	EnvironmentID uint        `gorm:"foreignKey:EnvironmentID" json:"environment_id"`
	Environment   Environment `json:"environment"`
	Name          string      `gorm:"type:varchar(100);not null" json:"name"`
	Value         string      `gorm:"type:varchar(255)" json:"value"`
	Description   string      `gorm:"type:varchar(255)" json:"description"`
}
