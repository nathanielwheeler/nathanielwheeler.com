package controllers

import (
	"fmt"
	"net/http"

	"nathanielwheeler.com/views"
)

// Subscribers :
type Subscribers struct {
	NewView *views.View
}

// SubscribeForm :
type SubscribeForm struct {
	Email string `schema:"email"`
}

// New : GET /subscribe
// — Renders a new form view for a potential subscriber
func (sub *Subscribers) New(res http.ResponseWriter, req *http.Request) {
	if err := sub.NewView.Render(res, nil); err != nil {
		// TODO don't panic && give feedback to subscriber
		panic(err)
	}
}

// Create : POST /subscribe
// — Used to process the subscription form when a user tries to subscribe
func (sub *Subscribers) Create(res http.ResponseWriter, req *http.Request) {
	var form SubscribeForm
	if err := parseForm(req, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(res, "Email is", form.Email)
}

// NewSubscribers : Initializes the view for subscribers
func NewSubscribers() *Subscribers {
	return &Subscribers{
		NewView: views.NewView("app", "subscribers/new"),
	}
}
