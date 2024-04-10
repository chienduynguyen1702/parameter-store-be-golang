package models

import (
	"gorm.io/gorm"
)

type Agent struct {
	gorm.Model
	ProjectID       uint
	Name            string `gorm:"type:varchar(100);not null"`
	APIToken        string `gorm:"type:varchar(100)"`
	StageID         uint   `gorm:"foreignKey:StageID;not null"`
	Stage           Stage
	EnvironmentID   uint `gorm:"foreignKey:EnvironmentID;not null"`
	Environment     Environment
	CICDWorflowName string `gorm:"type:varchar(100);not null"`
}
