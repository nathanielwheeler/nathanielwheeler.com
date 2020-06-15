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

// #region SERVICE

// PostsService will handle business rules for posts.
type PostsService interface {
	PostsDB
}

type postsService struct {
	PostsDB
}

// NewPostsService is 
func NewPostsService(db *gorm.DB) PostsService {
	return &postsService{
		PostsDB: &postsValidator{
			PostsDB: &postsGorm{
				db: db,
			},
		},
	}
}

// #endregion

// #region GORM

//		#region GORM CONFIG

// PostsDB will handle database interaction for posts.
type PostsDB interface {
	Create(post *Post) error
}

type postsGorm struct {
	db *gorm.DB
}

// Ensure that postsGorm always implements PostsDB interface
var _ PostsDB = &postsGorm{}

// #endregion

//		#region GORM METHODS

func (pg *postsGorm) Create(post *Post) error {
	return pg.db.Create(post).Error
}

// #endregion

// #endregion

// #region VALIDATOR

type postsValidator struct {
	PostsDB
}

// #endregion
