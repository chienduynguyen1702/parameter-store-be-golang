package models

import (
	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Password  string
	Address   string
	Posts     []Post `gorm:"many2many:author_posts;"`
}
