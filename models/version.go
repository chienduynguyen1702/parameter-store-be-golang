package models

import "gorm.io/gorm"

type Version struct {
	gorm.Model
	Number      string      `gorm:"type:varchar(100);not null" json:"number"`
	Name        string      `gorm:"index:agent_name_project_id;type:varchar(100);not null" json:"name"`
	// ProjectID   uint        `gorm:"index:agent_name_project_id;foreignKey:ProjectID" json:"project_id"`
	Description string      `gorm:"type:text" json:"description"`
	Parameters  []Parameter `gorm:"many2many:version_parameters" json:"parameters"`
	ProjectID   uint        `json:"project_id"`
}
