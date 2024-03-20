package models

import (
	"time"

	"gorm.io/gorm"
)

// Project represents the projects table
type Project struct {
	gorm.Model
	OrganizationID uint
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	Name           string       `gorm:"type:varchar(100);not null"`
	StartAt        time.Time
	Description    string `gorm:"type:text"`
	CurrentSprint  int
	RepoURL        string        `gorm:"type:varchar(100);not null"`
	Versions       []Version     `gorm:"one2many:project_versions;"`
	Stages         []Stage       `gorm:"one2many:project_stages;"`
	Environments   []Environment `gorm:"one2many:project_environments;"`
	Agent          []Agent       `gorm:"one2many:project_agents;"`
}
