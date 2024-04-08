package models

import (
	"gorm.io/gorm"
)

// Token model
type Token struct {
	gorm.Model
	UserID uint   `gorm:"not null" json:"user_id"`
	Token  string `gorm:"type:varchar(255);not null" json:"token"`
	User   User   `gorm:"foreignKey:UserID" json:"user"`
}
