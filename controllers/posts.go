package controllers

import (
	"nathanielwheeler.com/models"
	"nathanielwheeler.com/views"
)

// Posts will hold information about views and services
type Posts struct {
	New *views.View
	ps  models.PostsService
}

// NewPosts is a constructor for Posts struct
func NewPosts(ps models.PostsService) *Posts {
	return &Posts{
		New: views.NewView("app", "posts/new"),
		ps:  ps,
	}
}
