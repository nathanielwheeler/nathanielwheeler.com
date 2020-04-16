package controllers

import (
	"fmt"
	"net/http"

	"nathanielwheeler.com/models"
	"nathanielwheeler.com/views"
)

// Subscribers :
type Subscribers struct {
	NewView *views.View
	ss      *models.SubsService
}

// SubscribeForm :
type SubscribeForm struct {
	Email       string `schema:"email"`
	EveryUpdate bool   `schema:"every-update"`
}

// New : GET /subscribe
// — Renders a new form view for a potential subscriber
func (s *Subscribers) New(res http.ResponseWriter, req *http.Request) {
	if err := s.NewView.Render(res, nil); err != nil {
		// TODO don't panic && give feedback to subscriber
		panic(err)
	}
}

// Create : POST /subscribe
// — Used to process the subscription form when a user tries to subscribe
func (s *Subscribers) Create(res http.ResponseWriter, req *http.Request) {
	var form SubscribeForm
	if err := parseForm(req, &form); err != nil {
		panic(err)
	}
	sub := models.Subscriber{
		Email:       form.Email,
	}
	if err := s.ss.Create(&sub); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(res, "Subscriber is", sub)
}

// NewSubscribers : Initializes the view for subscribers
func NewSubscribers(ss *models.SubsService) *Subscribers {
	return &Subscribers{
		NewView: views.NewView("app", "subscribers/new"),
		ss:      ss,
	}
}
