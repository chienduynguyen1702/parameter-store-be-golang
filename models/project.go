package models

import (
	"time"

	"gorm.io/gorm"
)

// Project represents the projects table
type Project struct {
	gorm.Model
	OrganizationID  uint              `gorm:"not null" json:"organization_id"`
	Name            string            `gorm:"type:varchar(100);not null" json:"name"`
	StartAt         time.Time         `json:"start_at"`
	Address         string            `gorm:"type:varchar(100)" json:"address"`
	Status          string            `gorm:"type:varchar(100)" json:"status"`
	Description     string            `gorm:"type:text" json:"description"`
	CurrentSprint   string            `gorm:"type:varchar(100)" json:"current_sprint"`
	RepoURL         string            `gorm:"type:varchar(100)" json:"repo_url"`
	RepoApiToken    string            `gorm:"type:varchar(100)" json:"repo_api_token"`
	IsArchived      bool              `gorm:"default:false" json:"is_archived"`
	ArchivedBy      string            `json:"archived_by"` // foreign key to user model
	ArchivedAt      time.Time         `gorm:"type:timestamp;" json:"archived_at"`
	LatestVersionID uint              `gorm:"default:null" json:"latest_version_id"`
	AutoUpdate      bool              `gorm:"default:true" json:"auto_update"`
	LatestVersion   Version           `gorm:"-" json:"latest_version"`
	Stages          []Stage           `gorm:"foreignKey:ProjectID" json:"stages"`
	Environments    []Environment     `gorm:"foreignKey:ProjectID" json:"environments"`
	Versions        []Version         `gorm:"foreignKey:ProjectID" json:"versions"`
	Agents          []Agent           `gorm:"foreignKey:ProjectID" json:"agents"`
	Parameters      []Parameter       `gorm:"foreignKey:ProjectID" json:"parameters"`
	UserRoles       []UserRoleProject `gorm:"foreignKey:ProjectID" json:"user_roles"`
	Logs            []ProjectLog      `gorm:"foreignKey:ProjectID" json:"logs"`
	Workflows       []Workflow        `gorm:"foreignKey:ProjectID" json:"workflows"`
}
