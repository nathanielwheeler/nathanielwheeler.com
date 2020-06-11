package controllers

import (
	"fmt"
	"net/http"

	"nathanielwheeler.com/models"
	"nathanielwheeler.com/rand"
	"nathanielwheeler.com/views"
)

// NewUsers initializes the view for users
func NewUsers(us *models.UserService) *Users {
	return &Users{
		RegisterView: views.NewView("app", "users/register"),
		LoginView:    views.NewView("app", "users/login"),
		us:           us,
	}
}

// Users holds reference for the Users view and service.
type Users struct {
	RegisterView *views.View
	LoginView    *views.View
	us           *models.UserService
}

// RegisterForm : GET /register
// — Renders a new registration form for a potential user
func (u *Users) RegisterForm(res http.ResponseWriter, req *http.Request) {
	if err := u.RegisterView.Render(res, nil); err != nil {
		// TODO don't panic && give feedback to user
		panic(err)
	}
}

// RegistrationForm is used to transform a webform into a registration request
type RegistrationForm struct {
	Email    string `schema:"email"`
	Name     string `schema:"name"`
	Password string `schema:"password"`
}

// Register : POST /register
// — Used to process the signup form when a user tries to create a new user account
func (u *Users) Register(res http.ResponseWriter, req *http.Request) {
	var form RegistrationForm
	if err := parseForm(req, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// remember token
	err := u.signIn(res, &user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(res, req, "/cookietest", http.StatusFound)
}

// Login : POST /login
// — Used to process the login form when a user tries to log in as an existing user
func (u *Users) Login(res http.ResponseWriter, req *http.Request) {
	var form LoginForm
	if err := parseForm(req, &form); err != nil {
		panic(err)
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
		case models.ErrInvalidPassword:
			fmt.Fprintln(res, "Invalid email and/or password.")
		default:
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = u.signIn(res, user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(res, req, "/cookietest", http.StatusFound)
}

// LoginForm is used to transform a webform into a login request
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// signIn is used to sign the given user in via cookies
func (u *Users) signIn(res http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:  "remember_token",
		Value: user.Remember,
	}
	http.SetCookie(res, &cookie)
	return nil
}

// CookieTest is used to display cookies set on the current user
func (u *Users) CookieTest(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("remember_token")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(res, user)
}
