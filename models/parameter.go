package models

import (
	"gorm.io/gorm"
)

type Parameter struct {
	gorm.Model
	StageID       uint
	Stage         Stage
	EnvironmentID uint
	Environment   Environment
	Name          string `gorm:"type:varchar(100);not null"`
	Value         string `gorm:"type:varchar(255)"`
	Description   string `gorm:"type:varchar(255)"`
}
