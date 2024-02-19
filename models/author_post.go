package models

import (
	"gorm.io/gorm"
)

type Author_Post struct {
	gorm.Model
	AuthorID uint
	PostID   uint
}

// Path:
