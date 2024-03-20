package models

import "time"

type Environment struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(100);not null"`
	ProjectID uint
	CreateAt  time.Time `gorm:"default:current_timestamp"`
	UpdateAt  time.Time `gorm:"default:current_timestamp"`
	DeleteAt  *time.Time
}
