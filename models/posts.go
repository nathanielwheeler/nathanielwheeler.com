package models

import (
	"github.com/jinzhu/gorm"
)

// Post will hold all of the information needed for a blog post.
type Post struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

// PostsService will handle business rules for posts.
type PostsService interface {
	PostsDB
}

// PostsDB will handle database interaction for posts.
 type PostsDB interface {
	 Create(post *Post) error
 }

 type postsGorm struct {
	 db *gorm.DB
 }

 func (pg *postsGorm) Create(post *Post) error {
	// TODO: implement
	return nil
 }