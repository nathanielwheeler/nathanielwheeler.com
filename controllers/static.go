package controllers

import (
	"nathanielwheeler.com/views"
)

// Static : A type that holds the views of the static pages
type Static struct {
	Home,
	Resume *views.View
}

// NewStatic : Returns the initialized views of the static pages.
func NewStatic() *Static {
	return &Static{
		Home:   views.NewView("app", "views/static/home.html"),
		Resume: views.NewView("app", "views/static/resume.html"),
	}
}
