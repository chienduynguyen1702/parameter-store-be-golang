package model

import "gorm.io/gorm"

type Version struct {
	gorm.Model
	ProjectID   uint
	Description string      `gorm:"type:text"`
	Parameters  []Parameter `gorm:"many2many:version_parameters;foreignKey:ID;joinForeignKey:VersionID;References:ID;JoinReferences:ParameterID"`
}
