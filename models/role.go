package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string       `gorm:"type:varchar(100);not null"`
	Description string       `gorm:"type:text"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}
