package models

import (
	"time"

	"gorm.io/gorm"
)

// Project represents the projects table
type Project struct {
	gorm.Model
	OrganizationID  uint      `gorm:"foreignKey:OrganizationID;not null" json:"organization_id"`
	Name            string    `gorm:"type:varchar(100);not null" json:"name"`
	StartAt         time.Time `json:"start_at"`
	Address         string    `gorm:"type:varchar(100)" json:"address"`
	Status          string    `gorm:"type:varchar(100)" json:"status"`
	Description     string    `gorm:"type:text" json:"description"`
	CurrentSprint   string    `gorm:"type:varchar(100)" json:"current_sprint"`
	RepoURL         string    `gorm:"type:varchar(100)" json:"repo_url"`
	RepoApiToken    string    `gorm:"type:varchar(100)" json:"repo_api_token"`
	IsArchived      bool      `gorm:"default:false" json:"is_archived"`
	ArchivedBy      string    `gorm:"foreignKey:ArchivedBy" json:"archived_by"` // foreign key to user model
	ArchivedAt      time.Time `gorm:"type:timestamp;" json:"archived_at"`
	LatestVersionID uint      `gorm:"foreignKey:LatestVersionID" json:"latest_version"`
	LatestVersion   Version

	Stages       []Stage           `gorm:"many2many:project_stages;" json:"stages"`
	Environments []Environment     `gorm:"many2many:project_environments;" json:"environments"`
	Versions     []Version         `gorm:"one2many:project_versions;" json:"versions"`
	Agents       []Agent           `gorm:"one2many:project_agents;" json:"agents"`
	Parameters   []Parameter       `gorm:"one2many:project_parameters;" json:"parameters"`
	UserRoles    []UserRoleProject `gorm:"one2many:user_role_project;" json:"user_role"`
}
