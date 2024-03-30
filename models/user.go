package models

import (
	"gorm.io/gorm"
)

// User represents the users table
type User struct {
	gorm.Model
	Email               string `gorm:"type:varchar(255);unique;not null"`
	Username            string `gorm:"type:varchar(100);unique;not null"`
	Password            string `gorm:"type:varchar(255);not null"`
	Phone               string `gorm:"type:varchar(255);"`
	IsOrganizationAdmin bool   `gorm:"default:false"`                      // Assuming this field represents the ID of the organization the user is an admin of
	OrganizationID      uint   `gorm:"not null;foreignKey:OrganizationID"` // foreign key to organization model
}
