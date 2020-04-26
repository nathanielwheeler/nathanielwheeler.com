package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	// Not directly used, but needed to help gorm communicate with postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound : Indicates that a resource does not exist within postgres
	ErrNotFound = errors.New("models: resource not found")
	// ErrInvalidID : Returned when an invalid ID is provided to a method like Delete.
	ErrInvalidID = errors.New("models: ID provided was invalid")
)

// User : Model for people that want updates from my website and want to leave comments on my posts.
type User struct {
	gorm.Model
	Email string `gorm:"type:varchar(100);primary key"`
	Name  string
}

// UsersService : Processes the logic for users
type UsersService struct {
	db *gorm.DB
}

// NewUsersService : constructor for UsersService.  Initializes database connection
func NewUsersService(connectionStr string) (*UsersService, error) {
	db, err := gorm.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UsersService{
		db: db,
	}, nil
}

// #region SERVICE METHODS

// ByID : Gets a user given an ID.
func (us *UsersService) ByID(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Create : Creates the provided user and fills provided data fields
func (us *UsersService) Create(u *User) error {
	return us.db.Create(u).Error
}

// Update : Changes subscriber preferences
func (us *UsersService) Update(u *User) error {
	return us.db.Save(u).Error
}

// Delete : Removes the subscriber identified by the id
func (us *UsersService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	u := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&u).Error
}

// AutoMigrate : Attempts to automatically migrate the subscribers table
func (us *UsersService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// DestructiveReset : Destroys tables and calls AutoMigrate()
func (us *UsersService) DestructiveReset() error {
	err := us.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return us.AutoMigrate()
}

// Close : Shuts down the connection to database
func (us *UsersService) Close() error {
	return us.db.Close()
}

// #endregion
