package services

import (
	"github.com/jinzhu/gorm"
	// Since this is implicitly needed by gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Services will hold information about the varying services used in the services package.
type Services struct {
	User   UserService
	Posts  PostsService
	Images ImagesService
	db     *gorm.DB
}

// NewServices will accept a list of config functions to run.  Each function will accept a pointer to the current Services object, manipulate that object, returning an error if there is one.
func (s *server) NewServices(cfgs ...ServicesConfig) (error) {
	var services Services
	for _, cfg := range cfgs {
		if err := cfg(&services); err != nil {
			return nil, err
		}
	}
	s.services = &s
	return nil
}

// ServiceConfig is a type of functional option which returns an error.
type ServiceConfig func(*Services) error

// WithGorm is a functional option that will open a connection to GORM, returning an error if something goes wrong.
func WithGorm(dialect, connectionStr string) ServiceConfig {
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
func WithLogMode(mode bool) ServiceConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

// WithUser is a functional option that will construct a new user service, adding in pepper and HMAC key.
func WithUser(pepper, hmacKey string) ServiceConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper, hmacKey)
		return nil
	}
}

// WithPosts is a functional option that will construct a new posts service.
func WithPosts(isProd bool) ServiceConfig {
	return func(s *Services) error {
		s.Posts = NewPostsService(s.db, isProd)
		return nil
	}
}

// WithImages is a function option that will construct a new images service.
func WithImages() ServiceConfig {
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

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

func last(db *gorm.DB, dst interface{}) error {
	err := db.Last(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}