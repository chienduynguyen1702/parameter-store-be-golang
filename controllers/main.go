package controllers

import (
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

// SetDB sets the db object
func SetDB(database *gorm.DB) {
	DB = database
}
