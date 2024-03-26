package models

import (
	"gorm.io/gorm"
)

type Agent struct {
	gorm.Model
	ProjectID     uint
	Name          string `gorm:"type:varchar(100);not null"`
	APIToken      string `gorm:"type:varchar(100)"`
	StageID       uint
	Stage         Stage
	EnvironmentID uint
	Environment   Environment
}
