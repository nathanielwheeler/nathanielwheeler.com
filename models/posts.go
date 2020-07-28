package models

import (
	"github.com/jinzhu/gorm"
)

// #region ERRORS

/* TODO
- Need to make a private error type */

const (
	// ErrUserIDRequired indicates that there is a missing user ID
	ErrUserIDRequired modelError = "models: user ID is required"
	// ErrTitleRequired indicates that there is a missing title
	ErrTitleRequired modelError = "models: title is required"
)

// #endregion

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
	ByID(id uint) (*Post, error)
	ByYearAndTitle(year int, title string) (*Post, error)
	Create(post *Post) error
}

type postsGorm struct {
	db *gorm.DB
}

// Ensure that postsGorm always implements PostsDB interface
var _ PostsDB = &postsGorm{}

//		#endregion

//		#region GORM METHODS

// ByID will search the posts database for input id.
func (pg *postsGorm) ByID(id uint) (*Post, error) {
	var post Post
	db := pg.db.Where("id = ?", id)
	err := first(db, &post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// ByYearAndTitle will search the posts database for input year and title.
func (pg *postsGorm) ByYearAndTitle(year int, title string) (*Post, error) {
	var post Post
	db := pg.db.Where("title = ? AND created_at > '?-01-01 00:00:00-07' AND created_at <= '?-12-31 23:59:59-07'", title, year, year)
	err := first(db, &post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (pg *postsGorm) Create(post *Post) error {
	return pg.db.Create(post).Error
}

//		#endregion

// #endregion

// #region VALIDATOR

type postsValidator struct {
	PostsDB
}

//		#region DB VALIDATORS

func (pv *postsValidator) Create(post *Post) error {
	err := runPostsValFns(post,
		pv.userIDRequired,
		pv.titleRequired)
	if err != nil {
		return err
	}
	return pv.PostsDB.Create(post)
}

//		#endregion

//		#region VAL METHODS

type postsValFn func(*Post) error

func runPostsValFns(post *Post, fns ...postsValFn) error {
	for _, fn := range fns {
		if err := fn(post); err != nil {
			return err
		}
	}
	return nil
}

func (pv *postsValidator) userIDRequired(p *Post) error {
	if p.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (pv *postsValidator) titleRequired(p *Post) error {
	if p.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

//		#endregion

// #endregion
