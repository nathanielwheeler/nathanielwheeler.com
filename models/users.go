package models

import (
	"errors"

	"nathanielwheeler.com/hash"
	"nathanielwheeler.com/rand"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Not directly used, but needed to help gorm communicate with postgres
	"golang.org/x/crypto/bcrypt"
)

// TODO: randomize and store securely
const (
	hmacSecretKey = "secret-hmac-key"
	userPwPepper  = "secret-string"
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

// User is a model of the people using my website.
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"type:varchar(100);primary key"`
	Password     string `gorm:"-"` // Ensures that it won't be saved to database
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// UsersDB is used to interact with the users database
type UsersDB interface {
	// Methods for querying single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Used to close a DB connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

// UsersService processes the logic for users
type UsersService struct {
	UsersDB
}

type usersValidator struct {
	UsersDB
}

type usersGorm struct {
	db   *gorm.DB
	hmac hash.HMAC
}

// NewUsersService is a constructor for UsersService.  It embeds UsersDB (usersGorm) to handle database connections.
func NewUsersService(connectionStr string) (*UsersService, error) {
	ug, err := newUsersGorm(connectionStr)
	if err != nil {
		return nil, err
	}
	return &UsersService{
		UsersDB: ug,
	}, nil
}

// newUsersGorm is a constructor for UsersGorm.  Initializes database connection.
func newUsersGorm(connectionStr string) (*usersGorm, error) {
	db, err := gorm.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	return &usersGorm{
		db:   db,
		hmac: hmac,
	}, nil
}

// #region METHODS: UsersValidation

/* 
func (uv *usersValidator) ByID (id uint) (*User, error) {
	// Validate the ID
	if id <= 0 {
		return nil, errors.New("Invalid ID")
	}
	// If valid, call the next method in the chain and return its results.
	return uv.UsersDB.ByID(id)
} */

// #endregion

// #region METHODS: UsersDB

// ByID gets a user given an ID.
func (ug *usersGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByEmail gets a user given an email string
func (ug *usersGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up a user with the given remember token and returns that user.
func (ug *usersGorm) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := ug.hmac.Hash(token)
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create makes the provided user and fills provided data fields
func (ug *usersGorm) Create(user *User) error {
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
	user.RememberHash = ug.hmac.Hash(user.Remember)

	return ug.db.Create(user).Error
}

// Update changes subscriber preferences
func (ug *usersGorm) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = ug.hmac.Hash(user.Remember)
	}
	return ug.db.Save(user).Error
}

// Delete removes the subscriber identified by the id
func (ug *usersGorm) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// Close shuts down the connection to database
func (ug *usersGorm) Close() error {
	return ug.db.Close()
}

// AutoMigrate attempts to automatically migrate the subscribers table
func (ug *usersGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// DestructiveReset destroys tables and calls AutoMigrate()
func (ug *usersGorm) DestructiveReset() error {
	err := ug.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// #endregion



// Authenticate is used to authenticate a user with a provided email address and password.  Returns a user and an error message.
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

// #region HELPERS

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// #endregion
