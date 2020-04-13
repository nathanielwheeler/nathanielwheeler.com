package controllers

import (
	"net/http"
	"nathanielwheeler.com/views"
)

// Subscribers :
type Subscribers struct {
	NewView *views.View
}

// New : 
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