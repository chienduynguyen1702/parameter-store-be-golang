package models

import "gorm.io/gorm"

type Version struct {
	gorm.Model
	Number      string `gorm:"type:varchar(100);not null"`
	Name        string `gorm:"type:varchar(100);not null"`
	ProjectID   uint
	Description string      `gorm:"type:text"`
	Parameters  []Parameter `gorm:"many2many:version_parameters;foreignKey:ID;joinForeignKey:VersionID;References:ID;JoinReferences:ParameterID"`
}
