package models

import (
	"errors"

	"nathanielwheeler.com/hash"
	"nathanielwheeler.com/rand"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Not directly used, but needed to help gorm communicate with postgres
	"golang.org/x/crypto/bcrypt"
)

// #region ERRORS
var (
	// ErrNotFound : Indicates that a resource does not exist within postgres
	ErrNotFound = errors.New("models: resource not found")
	// ErrInvalidID : Returned when an invalid ID is provided to a method like Delete.
	ErrInvalidID = errors.New("models: ID provided was invalid")
	// ErrInvalidPassword : Returned when an invalid password is is used when attempting to authenticate a user.
	ErrInvalidPassword = errors.New("models: incorrect password")
)

// #endregion

// TODO: remove obvious pepper
var userPwPepper = "secret-string"

// TODO: remove obvious hmac key
var hmacSecretKey = "secret-hmac-key"

// User : Model for people that want updates from my website and want to leave comments on my posts.
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"type:varchar(100);primary key"`
	Password     string `gorm:"-"` // Ensures that it won't be saved to database
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null; unique_index"`
}

// UsersService : Processes the logic for users
type UsersService struct {
	db   *gorm.DB
	hmac hash.HMAC
}

// NewUsersService : constructor for UsersService.  Initializes database connection
func NewUsersService(connectionStr string) (*UsersService, error) {
	db, err := gorm.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	hmac := hash.NewHMAC(hmacSecretKey)

	return &UsersService{
		db:   db,
		hmac: hmac,
	}, nil
}

// #region SERVICE METHODS

// ByID : Gets a user given an ID.
func (us *UsersService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByEmail : Get a user given an email string
func (us *UsersService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// Create : Creates the provided user and fills provided data fields
func (us *UsersService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(
		pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}

	return us.db.Create(user).Error
}

// Update : Changes subscriber preferences
func (us *UsersService) Update(user *User) error {
	return us.db.Save(user).Error
}

// Delete : Removes the subscriber identified by the id
func (us *UsersService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

// Close : Shuts down the connection to database
func (us *UsersService) Close() error {
	return us.db.Close()
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

// Authenticate : Used to authenticate a user with a provided email address and password.  Returns a user and an error message.
func (us *UsersService) Authenticate(email, password string) (*User, error) {
	// Check if email exists
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	// If email found, compare password hashes. Return user or error statement.
	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash),
		[]byte(password+userPwPepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrInvalidPassword
	default:
		return nil, err
	}
}

// #endregion

// #region HELPERS

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// #endregion
