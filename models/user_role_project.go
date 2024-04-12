package models

import (
	"gorm.io/gorm"
)

// UserProjectRole model
type UserRoleProject struct {
	gorm.Model
	UserID    uint `gorm:"foreignKey:UserID,not null" json:"user_id"`
	ProjectID uint `gorm:"foreignKey:ProjectID,not null" json:"project_id"`
	RoleID    uint `gorm:"foreignKey:RoleID,not null" json:"role_id"`
	User      User `gorm:"foreignKey:UserID" json:"user"`
	Role      Role `gorm:"foreignKey:RoleID" json:"role"`
}
