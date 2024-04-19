package models

import (
	"time"

	"gorm.io/gorm"
)

type Environment struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Color       string    `gorm:"type:varchar(100)" json:"color"`
	ProjectID   uint      `gorm:"foreignKey:ProjectID" json:"project_id"`
	IsArchived  bool      `gorm:"default:false" json:"is_archived"`
	ArchivedBy  string    `gorm:"foreignKey:ArchivedBy" json:"archived_by"`
	ArchivedAt  time.Time `gorm:"type:timestamp;" json:"archived_at"`
}
