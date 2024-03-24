package model

type VersionParameter struct {
	ID          uint `gorm:"type:serial;primaryKey"`
	VersionID   uint `gorm:"type:serial;index"`
	ParameterID uint `gorm:"type:serial;index"`
}
