package models

import "time"

type Agent struct {
	ID            uint `gorm:"primaryKey"`
	ProjectID     uint
	Name          string    `gorm:"type:varchar(100);not null"`
	APIToken      string    `gorm:"type:varchar(100)"`
	CreateAt      time.Time `gorm:"default:current_timestamp"`
	UpdateAt      time.Time `gorm:"default:current_timestamp"`
	DeleteAt      *time.Time
	StageID       uint
	Stage         Stage
	EnvironmentID uint
	Environment   Environment
}
