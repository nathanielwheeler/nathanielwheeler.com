package controllers

import (
	"fmt"
	"net/http"

	"nathanielwheeler.com/models"
	"nathanielwheeler.com/views"
)

// NewUsers : Initializes the view for users
func NewUsers(us *models.UsersService) *Users {
	return &Users{
		NewView: views.NewView("app", "users/new"),
		LoginView: views.NewView("app", "users/login"),
		us:      us,
	}
}

// Users : Holds reference for the Users view and service.
type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UsersService
}

// #region FORMS

// SignupForm : This form is used to transform a webform into a registration request
type SignupForm struct {
	Email    string `schema:"email"`
	Name     string `schema:"name"`
	Password string `schema:"password"`
}

// LoginForm : This form is used to transform a webform into a login request
type LoginForm struct {
	Email string `schema:"email"`
	Password string `schema:"password"`
}

// #endregion

// New : GET /signup
// — Renders a new form view for a potential user
func (u *Users) New(res http.ResponseWriter, req *http.Request) {
	if err := u.NewView.Render(res, nil); err != nil {
		// TODO don't panic && give feedback to user
		panic(err)
	}
}

// Login : POST /login
// — Used to process the login form when a user tries to log in as an existing user
func (u *Users) Login(res http.ResponseWriter, req *http.Request) {
	form := LoginForm{}
	if err := parseForm(req, &form); err != nil {
		panic(err)
	}
}

// Create : POST /signup
// — Used to process the signup form when a user tries to create a new user account
func (u *Users) Create(res http.ResponseWriter, req *http.Request) {
	var form SignupForm
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
	fmt.Fprintln(res, "User is", user)
}