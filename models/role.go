package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string       `gorm:"type:varchar(100);not null" json:"name" binding:"required"`
	Description string       `gorm:"type:text" json:"description" binding:"required"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions" binding:"required"`
}
