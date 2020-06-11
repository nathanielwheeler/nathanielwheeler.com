package controllers

import (
	"fmt"
	"net/http"

	"nathanielwheeler.com/models"
	"nathanielwheeler.com/views"
)

// NewUsers : Initializes the view for users
func NewUsers(us *models.UserService) *Users {
	return &Users{
		RegisterView: views.NewView("app", "users/register"),
		LoginView:    views.NewView("app", "users/login"),
		us:           us,
	}
}

// Users : Holds reference for the Users view and service.
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

// RegistrationForm : This form is used to transform a webform into a registration request
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
	fmt.Fprintln(res, "User is", user)
}

// Login : POST /login
// — Used to process the login form when a user tries to log in as an existing user
func (u *Users) Login(res http.ResponseWriter, req *http.Request) {
	var form LoginForm
	if err := parseForm(req, &form); err != nil {
		panic(err)
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	switch err {
	// TODO Remove this error message
	case models.ErrNotFound:
		fmt.Fprintln(res, "Email not found")
	case models.ErrInvalidPassword:
		fmt.Fprintln(res, "Email and password do not match.")
	case nil:
		fmt.Fprintln(res, "User is", user)
	default:
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// LoginForm : This form is used to transform a webform into a login request
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// CookieTest is used to display cookies set on the current user
func (u *Users) CookieTest(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("email")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(res, "Email is:", cookie.Value)
	fmt.Fprintf(res, "%+v", cookie)
}
