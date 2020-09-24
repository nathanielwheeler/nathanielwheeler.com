package models

import (
	"github.com/jinzhu/gorm"
	// Since this is implicitly needed by gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Services will hold information about the varying services used in the models package.
type Services struct {
	User   UserService
	Posts  PostsService
	Images ImagesService
	db     *gorm.DB
}

// NewServices will accept a list of config functions to run.  Each function will accept a pointer to the current Services object, manipulate that object, returning an error if there is one.
func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

// ServicesConfig is a type of functional option which returns an error.
type ServicesConfig func(*Services) error

// WithGorm is a functional option that will open a connection to GORM, returning an error if something goes wrong.
func WithGorm(dialect, connectionStr string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionStr)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

// WithLogMode is a functional option that configure log mode with the database.
func WithLogMode(mode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

// WithUser is a functional option that will construct a new user service, adding in pepper and HMAC key.
func WithUser(pepper, hmacKey string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper, hmacKey)
		return nil
	}
}

// WithPosts is a functional option that will construct a new posts service.
func WithPosts() ServicesConfig {
	return func(s *Services) error {
		s.Posts = NewPostsService(s.db)
		return nil
	}
}

// WithImages is a function option that will construct a new images service.
func WithImages() ServicesConfig {
	return func(s *Services) error {
		s.Images = NewImagesService()
		return nil
	}
}

// Close shuts down the connection to the database
func (s *Services) Close() error {
	return s.db.Close()
}

// AutoMigrate will attempt to automatically migrate tables
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Post{}).Error
}

// DestructiveReset will drop tables and call AutoMigrate
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Post{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}
