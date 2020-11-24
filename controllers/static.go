package controllers

import (
	"nathanielwheeler.com/views"
)

// Static : A type that holds the views of the static pages
type Static struct {
	Resume,
	PrototypeThemeSystem *views.View
}

// NewStatic : Returns the initialized views of the static pages.
func NewStatic() *Static {
	return &Static{
		Resume:               views.NewView("app", "static/resume"),
		PrototypeThemeSystem: views.NewView("app", "static/prototypes/theme-system"),
	}
}
