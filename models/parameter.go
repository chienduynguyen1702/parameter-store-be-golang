package models

import (
	"gorm.io/gorm"
)

type Parameter struct {
	gorm.Model
	StageID       uint   `gorm:"foreignKey:StageID" json:"stage_id"`
	EnvironmentID uint   `gorm:"foreignKey:EnvironmentID" json:"environment_id"`
	Name          string `gorm:"type:varchar(100);not null" json:"name"`
	Value         string `gorm:"type:varchar(255)" json:"value"`
	Description   string `gorm:"type:varchar(255)" json:"description"`
	ProjectID     uint   `gorm:"foreignKey:ProjectID" json:"project_id"`
}
