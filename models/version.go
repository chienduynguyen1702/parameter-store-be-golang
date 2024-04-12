package models

import "gorm.io/gorm"

type Version struct {
	gorm.Model
	Number      string      `gorm:"type:varchar(100);not null" json:"number"`
	Name        string      `gorm:"type:varchar(100);not null" json:"name"`
	ProjectID   uint        `gorm:"foreignKey:ProjectID" json:"project_id"`
	Description string      `gorm:"type:text" json:"description"`
	Parameters  []Parameter `gorm:"many2many:version_parameters" json:"parameters"`
}
