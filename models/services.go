package models

import (
	"github.com/jinzhu/gorm"
)

// Services will hold information about the varying services used in the models package.
type Services struct {
	Posts  PostsService
	User   UserService
	Images ImageService
	db     *gorm.DB
}

// NewServices is a constructor for services.
func NewServices(connectionStr string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User:   NewUserService(db),
		Posts:  NewPostsService(db),
		Images: NewImageService(),
		db:     db,
	}, nil
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
