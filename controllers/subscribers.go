package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"

	"nathanielwheeler.com/views"
)

// Subscribers :
type Subscribers struct {
	NewView *views.View
}

// SubscribeForm :
type SubscribeForm struct {
	Email string `schema: "email"`
}

// New : GET /subscribe
// — Renders a new form view for a potential subscriber
func (sub *Subscribers) New(res http.ResponseWriter, req *http.Request) {
	if err := sub.NewView.Render(res, nil); err != nil {
		// TODO don't panic && give feedback to subscriber
		panic(err)
	}

	decoder := schema.NewDecoder()
	form := SubscribeForm{}
	if err := decoder.Decode(&form, req.PostForm); err != nil {
		panic(err)
	}
	fmt.Fprintln(res, form)
}

// Create : POST /subscribe
// — Used to process the subscription form when a user tries to subscribe
func (sub *Subscribers) Create(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		panic(err)
	}
	fmt.Fprintln(res, req.PostForm["email"])
}

// NewSubscribers : Initializes the view for subscribers
func NewSubscribers() *Subscribers {
	return &Subscribers{
		NewView: views.NewView("app", "views/subscribers/new.html"),
	}
}
