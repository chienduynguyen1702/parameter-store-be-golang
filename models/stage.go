package models

import (
	"gorm.io/gorm"
)

type Stage struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:text"`
	Color       string `gorm:"type:varchar(100)"`
	ProjectID   uint
}
