package controllers

import "gorm.io/gorm"

var (
	db *gorm.DB
)

// SetDB sets the db object
func SetDB(database *gorm.DB) {
	db = database
}
