package model

import (
	"time"

	"gorm.io/gorm"
)

// Organization model
type Organization struct {
	gorm.Model
	Name              string `gorm:"type:varchar(100);not null"`
	AliasName         string `gorm:"type:varchar(100)"`
	EstablishmentDate time.Time
	Description       string    `gorm:"type:text"`
	Projects          []Project `gorm:"one2many:organization_prs;"`
}
