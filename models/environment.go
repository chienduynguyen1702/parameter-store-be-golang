package models

import (
	"gorm.io/gorm"
)

type Environment struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Color       string `gorm:"type:varchar(100)" json:"color"`
}
