package controllers

import (
	"fmt"
	"net/http"

	"nathanielwheeler.com/models"
	"nathanielwheeler.com/views"
)

// Users : Holds reference for the Users view and service.
type Users struct {
	NewView *views.View
	us      *models.UsersService
}

// SignupForm : This form is used to transform a webform into something my server can use
type SignupForm struct {
	Email    string `schema:"email"`
	Name     string `schema:"name"`
	Password string `schema:"password"`
}

// New : GET /signup
// — Renders a new form view for a potential user
func (u *Users) New(res http.ResponseWriter, req *http.Request) {
	if err := u.NewView.Render(res, nil); err != nil {
		// TODO don't panic && give feedback to user
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

// NewUsers : Initializes the view for users
func NewUsers(us *models.UsersService) *Users {
	return &Users{
		NewView: views.NewView("app", "users/new"),
		us:      us,
	}
}
