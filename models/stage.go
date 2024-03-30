package models

import (
	"gorm.io/gorm"
)

type Stage struct {
	gorm.Model
	Name      string `gorm:"type:varchar(100);not null"`
	ProjectID uint
}
