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

// Create : POST /subscribe
// — Used to process the subscription form when a user tries to subscribe
func (sub *Subscribers) Create(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "test response")
}

// New : Renders a new form view for a potential subscriber
func (sub *Subscribers) New(res http.ResponseWriter, req *http.Request) {
	err := sub.NewView.Render(res, nil)
	if err != nil {
		// TODO don't panic && give feedback to subscriber
		panic(err)
	}
}

// NewSubscribers : Initializes the view for subscribers
func NewSubscribers() *Subscribers {
	return &Subscribers {
		NewView: views.NewView("app", "views/subscribers/new.html"),
	}
}