package models

import "time"

type Organization struct {
	ID                uint      `gorm:"primaryKey"`
	Name              string    `gorm:"type:varchar(100);not null"`
	AliasName         string    `gorm:"type:varchar(100)"`
	CreateAt          time.Time `gorm:"default:current_timestamp"`
	DeleteAt          *time.Time
	UpdateAt          time.Time `gorm:"default:current_timestamp"`
	EstablishmentDate time.Time
	Description       string    `gorm:"type:text"`
	Projects          []Project `gorm:"one2many:organization_prs;"`
}
