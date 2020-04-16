package models

import (
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	PostID uint `gorm:"not_null;index"`
	Title string `gorm:"not_null"`
}

type PostService interface {
	PostDB
}

type PostDB interface {
	Create(post *Post) error
}

type postGorm struct {
	db *gorm.DB
}