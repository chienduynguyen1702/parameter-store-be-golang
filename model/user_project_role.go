package model

import (
	"gorm.io/gorm"
)

// UserProjectRole model
type UserProjectRole struct {
	gorm.Model
	UserID    uint    `gorm:"not null"`
	ProjectID uint    `gorm:"not null"`
	RoleID    string  `gorm:"not null"`
	User      User    `gorm:"foreignKey:UserID"`
	Project   Project `gorm:"foreignKey:ProjectID"`
	Role      Role    `gorm:"foreignKey:RoleID"`
}
