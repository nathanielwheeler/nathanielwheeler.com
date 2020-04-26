package controllers

import (
	"fmt"
	"net/http"

	"nathanielwheeler.com/models"
	"nathanielwheeler.com/views"
)

// Users :
type Users struct {
	NewView *views.View
	us      *models.UsersService
}

// FIXME Web form must be updated

// SignupForm :
type SignupForm struct {
	Email       string `schema:"email"`
	EveryUpdate bool   `schema:"every-update"`
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
// — Used to process the subscription form when a user tries to signup
func (u *Users) Create(res http.ResponseWriter, req *http.Request) {
	var form SubscribeForm
	if err := parseForm(req, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Email:       form.Email,
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
