// models/author_post.go

package models

// "gorm.io/gorm"

type AuthorPost struct {
	// gorm.Model
	AuthorID uint
	PostID   uint
}
