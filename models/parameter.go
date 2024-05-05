package models

import (
	"time"

	"gorm.io/gorm"
)

type Parameter struct {
	gorm.Model
	StageID       uint      `gorm:"foreignKey:StageID" json:"stage_id"`
	EnvironmentID uint      `gorm:"foreignKey:EnvironmentID" json:"environment_id"`
	Name          string    `gorm:"type:varchar(100);not null" json:"name"`
	Value         string    `gorm:"type:varchar(255)" json:"value"`
	Description   string    `gorm:"type:varchar(255)" json:"description"`
	ProjectID     uint      `gorm:"foreignKey:ProjectID" json:"project_id"`
	IsArchived    bool      `gorm:"default:false" json:"is_archived"`
	ArchivedBy    string    `gorm:"foreignKey:ArchivedBy" json:"archived_by"` // foreign key to user model
	ArchivedAt    time.Time `gorm:"type:timestamp;" json:"archived_at"`
	IsApplied     bool      `gorm:"default:false" json:"is_applied"`

	// UpdatedBy   User		`gorm:"foreignKey:UpdatedBy" json:"updated_by"` // foreign key to user model
	Stage       Stage       `gorm:"foreignKey:StageID" json:"stage"`
	Environment Environment `gorm:"foreignKey:EnvironmentID" json:"environment"`
}
