package models

import (
	"gorm.io/gorm"
)

// Token model
type Token struct {
	gorm.Model
	UserID uint   `gorm:"not null"`
	Token  string `gorm:"type:varchar(255);not null"`
	User   User   `gorm:"foreignKey:UserID"`
}
