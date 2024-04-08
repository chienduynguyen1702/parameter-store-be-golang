package models

import (
	"gorm.io/gorm"
)

type Stage struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Color       string `gorm:"type:varchar(100)" json:"color"`
	ProjectID   uint   `gorm:"foreignKey:ProjectID" json:"project_id"`
}
