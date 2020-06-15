package models

import (
	"github.com/jinzhu/gorm"
)

// Services will hold information about the varying services used in the models package.
type Services struct {
	Posts PostsService
	User  UserService
}

// NewServices is a constructor for services.
func NewServices(connectionStr string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User:  NewUserService(db),
		Posts: &postsGorm{},
	}, nil
}
