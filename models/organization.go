package models

import (
	"time"

	"gorm.io/gorm"
)

// Organization model
type Organization struct {
	gorm.Model
	Name              string    `gorm:"type:varchar(100);not null" json:"name"`
	AliasName         string    `gorm:"type:varchar(100)" json:"alias_name"`
	EstablishmentDate time.Time `json:"establishment_date"`
	Description       string    `gorm:"type:text" json:"description"`
	Address           string    `gorm:"type:text" json:"address"`

	Projects []Project `gorm:"one2many:organization_prs;" json:"projects"`
	Users    []User    `gorm:"one2many:organization_usr;" json:"users"`
}
