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

// UserDB is used to interact with the users database.
type UserDB interface {
	// methods for single user queries
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)
	// methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
	// methods for DB management
	Close() error
	AutoMigrate() error
	DestructiveReset() error
}

// #region SERVICE

// UserService is a set of methods used to handle business rules of the user model
type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

// userService processes business rules for users
type userService struct {
	UserDB
}

// NewUserService : constructor for userService.  Initializes database connection
func NewUserService(connectionStr string) (UserService, error) {
	ug, err := newUserGorm(connectionStr)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := &userValidator{
		hmac: hmac,
		UserDB: ug,
	}
	return &userService{
		UserDB: uv,
	}, nil
}

// Authenticate : Used to authenticate a user with a provided email address and password.  Returns a user and an error message.
func (us *userService) Authenticate(email, password string) (*User, error) {
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

// #region GORM

// userGORM represents the database interaction layer
type userGorm struct {
	db   *gorm.DB
}

func newUserGorm(connectionStr string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &userGorm{
		db:   db,
	}, nil
}

// ByID gets a user given an ID.
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByEmail gets a user given an email string
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember gets a user given a remember token hash
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create takes in a validated user and adds it the database
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Update will update details of the user
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Delete removes the user identified by the id
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// Close shuts down the connection to database
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// AutoMigrate : Attempts to automatically migrate the subscribers table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// DestructiveReset : Destroys tables and calls AutoMigrate()
func (ug *userGorm) DestructiveReset() error {
	err := ug.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return ug.AutoMigrate()
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// #endregion

// #region VALIDATION

// userValidator represents the validation layer.  It also handles normalization.
type userValidator struct {
	UserDB
	hmac hash.HMAC
}

// ByRemember will hash the remember token then call ByRemember in the UserDB layer.
func (uv *userValidator) ByRemember(token string) (*User, error) {
	rememberHash := uv.hmac.Hash(token)
	return uv.UserDB.ByRemember(rememberHash)
}

// Create will make the provided user and backfill data like the ID, CreatedAt, and UpdatedAt fields.
func (uv *userValidator) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
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
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return uv.UserDB.Create(user)
}

// Update will hash a remember token if one is provided attached to the user object
func (uv *userValidator) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = uv.hmac.Hash(user.Remember)
	}
	return uv.UserDB.Update(user)
}

// Delete will delete the user with the provided ID
func (uv *userValidator) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	return uv.UserDB.Delete(id)
}

// #endregion
