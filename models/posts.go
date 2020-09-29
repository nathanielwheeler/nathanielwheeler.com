package models

import (
	"bytes"
	"io/ioutil"

	"github.com/jinzhu/gorm"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	_ "github.com/yuin/goldmark/extension" // Needed for goldmark extensions
	"github.com/yuin/goldmark/parser"
)

// Post will hold all of the information needed for a blog post.
type Post struct {
	gorm.Model
	Title       string                 `gorm:"not_null"`
	URLPath     string                 `gorm:"not_null"`
	FilePath    string                 `gorm:"not_null"`
	Body        string                 `gorm:"-"` // Not stored in database
	MetaData    map[string]interface{} `gorm:"-"`
}

// #region SERVICE

// PostsService will handle business rules for posts.
type PostsService interface {
  PostsDB
  ParseMD(*Post) error
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

// ParseMD will parse the associated markdown of a post.
func (ps *postsService) ParseMD(post *Post) error {
	data, err := ioutil.ReadFile(post.FilePath)
	if err != nil {
		return err
	}
	md := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	var buf bytes.Buffer
	ctx := parser.NewContext()
	if err := md.Convert([]byte(data), &buf, parser.WithContext(ctx)); err != nil {
		return err
	}

	post.Body = buf.String()
	post.MetaData = meta.Get(ctx)

	return nil
}

// #endregion

// #region GORM

//    #region GORM CONFIG

// PostsDB will handle database interaction for posts.
type PostsDB interface {
	ByID(id uint) (*Post, error)
	ByURL(urlpath string) (*Post, error)
	ByLatest() (*Post, error)
	GetAll() ([]Post, error)
	Create(post *Post) error
	Update(post *Post) error
	Delete(id uint) error
}

type postsGorm struct {
	db *gorm.DB
}

// Ensure that postsGorm always implements PostsDB interface
var _ PostsDB = &postsGorm{}

//    #endregion

//    #region GORM METHODS

// ByID will search the posts database for a post using input ID.
func (pg *postsGorm) ByID(id uint) (*Post, error) {
	var post Post
	db := pg.db.Where("id = ?", id)
	err := first(db, &post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// ByURL will search the posts database for input url string.
func (pg *postsGorm) ByURL(urlpath string) (*Post, error) {
	var post Post
	db := pg.db.Where("url_path = ?", urlpath)
	err := first(db, &post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// ByLatest will get the most recent post (by CreatedAt)
func (pg *postsGorm) ByLatest() (*Post, error) {
	var post Post
	pg.db.Raw(`SELECT * FROM "posts"  WHERE "posts"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT 1`).Scan(&post)
	return &post, nil
}

// GetAll will return all posts from newest to oldest.
func (pg *postsGorm) GetAll() ([]Post, error) {
	var posts []Post
	if err := pg.db.Order("created_at").Find(&posts).Error; err != nil {
		return nil, err
	}
	for i, v := 0, len(posts)-1; i < v; i, v = i+1, v-1 {
		posts[i], posts[v] = posts[v], posts[i]
	}
	return posts, nil
}

// Create will add a post to the database
func (pg *postsGorm) Create(post *Post) error {
	return pg.db.Create(post).Error
}

// Update will edit a post in a database
func (pg *postsGorm) Update(post *Post) error {
	return pg.db.Save(post).Error
}

// Delete will remove a post from default queries.
/* Really, it will add a timestamp for deleted_at, which will exclude the post from normal queries. */
func (pg *postsGorm) Delete(id uint) error {
	post := Post{Model: gorm.Model{ID: id}}
	return pg.db.Delete(&post).Error
}

//    #endregion

// #endregion

// #region VALIDATOR

type postsValidator struct {
	PostsDB
}

/*
VALIDATORS TODO
Ensure that title doesn't already exist in database (within year)
Ensure that title doesn't have any underscores in it
*/

//    #region DB VALIDATORS

func (pv *postsValidator) Create(post *Post) error {
	err := runPostsValFns(post,
		pv.titleRequired)
	if err != nil {
		return err
	}
	return pv.PostsDB.Create(post)
}

func (pv *postsValidator) Update(post *Post) error {
	err := runPostsValFns(post,
		pv.titleRequired)
	if err != nil {
		return err
	}
	return pv.PostsDB.Update(post)
}

func (pv *postsValidator) Delete(id uint) error {
	var post Post
	post.ID = id
	if err := runPostsValFns(&post, pv.nonZeroID); err != nil {
		return err
	}
	return pv.PostsDB.Delete(post.ID)
}

//    #endregion

//    #region VAL METHODS

type postsValFn func(*Post) error

func runPostsValFns(post *Post, fns ...postsValFn) error {
	for _, fn := range fns {
		if err := fn(post); err != nil {
			return err
		}
	}
	return nil
}

func (pv *postsValidator) titleRequired(p *Post) error {
	if p.Title == "" {
		return errTitleRequired
	}
	return nil
}

func (pv *postsValidator) nonZeroID(post *Post) error {
	if post.ID <= 0 {
		return errIDInvalid
	}
	return nil
}

//    #endregion

// #endregion
