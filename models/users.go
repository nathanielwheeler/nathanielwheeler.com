package models

import (
	"regexp"
	"strings"

	"nathanielwheeler.com/hash"
	"nathanielwheeler.com/rand"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User : Model for people that want updates from my website and want to leave comments on my posts.
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"type:varchar(100);primary key"`
	Password     string `gorm:"-"` // Ensures that it won't be saved to database
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null; unique_index"`
	IsAdmin      bool   `gorm:"default:false"`
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
	pepper string
}

// NewUserService : constructor for userService.  Calls constructors for user gorm and user validator.
func NewUserService(db *gorm.DB, pepper, hmacKey string) UserService {
	ug := &userGorm{db}
	hmac := hash.NewHMAC(hmacKey)
	uv := newUserValidator(ug, hmac, pepper)
	return &userService{
		UserDB: uv,
		pepper: pepper,
	}
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
		[]byte(password+us.pepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrPasswordInvalid
	default:
		return nil, err
	}
}

// #endregion

// #region GORM

// userGORM represents the database interaction layer
type userGorm struct {
	db *gorm.DB
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

// #endregion

// #region VALIDATION

// userValidator represents the validation layer.  It also handles normalization.
type userValidator struct {
	UserDB
	hmac       hash.HMAC
	pepper     string
	emailRegex *regexp.Regexp
	// TODO: make regex for password validation
	// passwordRegex *regexp.Regexp
}

// Constructor for userValidator layer.  Needed so that I can compile regex and assign it.
func newUserValidator(udb UserDB, hmac hash.HMAC, pepper string) *userValidator {
	return &userValidator{
		UserDB: udb,
		hmac:   hmac,
		pepper: pepper,
		emailRegex: regexp.MustCompile(
			`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
		// TODO figure out what I need for this
		/* Requirements:
		   - A-z
		   - 0-9
		   - !@#$%^&*()-_+={}[]\|:;'".,/<>?`~
		   - spaces
		*/
		/* passwordRegex: regexp.MustCompile(
		   ``), */
	}
}

// #region DB VALIDATORS

// ByEmail will normalize an email address before passing it to the database layer.
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	err := runUserValFns(&user, uv.normalizeEmail)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

// ByRemember will hash the remember token then call ByRemember in the UserDB layer.
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFns(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

// Create will make the provided user and backfill data like the ID, CreatedAt, and UpdatedAt fields.
func (uv *userValidator) Create(user *User) error {
	err := runUserValFns(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRememberIfUnset,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

// Update will hash a remember token if one is provided attached to the user object
func (uv *userValidator) Update(user *User) error {
	err := runUserValFns(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// Delete will delete the user with the provided ID, assuming that uint is greater than zero.
func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFns(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

// #endregion

// #region VAL METHODS

type userValFn func(*User) error

// runUserValFns will take a user object and loop through input validation functions to ensure that all data is proper.
func runUserValFns(user *User, fns ...userValFn) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

// bcryptPassword will hash a user's password with pepper and salt it with bcrypt.
func (uv *userValidator) bcryptPassword(user *User) error {
	// Will not run if password isn't provided in the provided user
	if user.Password == "" {
		return nil
	}

	pwBytes := []byte(user.Password + uv.pepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	return nil
}

// hmacRemember will hash a remember token if one is provided
func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

/* setRememberIfUnset will check the input user object to see if it has a remember token.
If it does, it returns nil.
If it doesn't, it calls rand.RememberToken() and adds the returned hash to the user. */
func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.RememberHash = token
	return nil
}

// idGreaterThan is to make sure that a uint of zero never reaches the delete function of Gorm, which would delete the entire database.
func (uv *userValidator) idGreaterThan(n uint) userValFn {
	return userValFn(func(user *User) error {
		if user.ID <= n {
			return errIDInvalid
		}
		return nil
	})
}

// normalizeEmail will ensure that email addresses are stored in a standardized format in the database.
func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

// requireEmail requires an email to be entered before continuing.
func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return errEmailRequired
	}
	return nil
}

// emailFormat checks if provided email address is in a proper form
func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		// This is a bit weird, but there may be instances where I accept an optional email field, but still need to make sure I have a proper email.
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return errEmailInvalid
	}
	return nil
}

// emailIsAvail checks if the email address entered is not already taken
func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	// This means the email address is available.  Continue.
	if err == ErrNotFound {
		return nil
	}
	// Other errors are bad.  Return error.
	if err != nil {
		return err
	}
	// Check if the existing email belongs to someone else.
	if user.ID != existing.ID {
		return errEmailTaken
	}
	return nil
}

// passwordMinLength enforces a minimum length on password. NOTE: It MUST be run before the password hasher.
func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return errPasswordTooShort
	}
	return nil
}

// passwordRequired returns an error if the password field passed in is empty
func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return errPasswordRequired
	}
	return nil
}

// passwordHashRequired returns if a hash for whatever reason wasn't made
func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return errPasswordRequired
	}
	return nil
}

// rememberMinBytes will check if remember token has 32 bytes
func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}
	if n < 32 {
		return errRememberTooShort
	}
	return nil
}

// rememberHashRequired checks if there is a remember token hash
func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return errRememberRequired
	}
	return nil
}

// #endregion

// #endregion
