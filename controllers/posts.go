package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"nathanielwheeler.com/context"
	"nathanielwheeler.com/models"
	"nathanielwheeler.com/views"

	"github.com/gorilla/mux"
)

// Posts will hold information about views and services
type Posts struct {
	New      *views.View
	ShowView *views.View
	ps       models.PostsService
}

// NewPosts is a constructor for Posts struct
func NewPosts(ps models.PostsService) *Posts {
	return &Posts{
		New:      views.NewView("app", "posts/new"),
		ShowView: views.NewView("app", "posts/show"),
		ps:       ps,
	}
}

// PostForm will hold information for creating a new post
type PostForm struct {
	Title string `schema:"title"`
}

// Show : GET /posts/:title
func (p *Posts) Show(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	titleVar := vars["title"]
	title := strings.Replace(titleVar, "-", " ", -1)

	post, err := p.ps.ByTitle(title)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(res, "Post not found", http.StatusNotFound)
		default:
			http.Error(res, "Something bad happened.", http.StatusInternalServerError)
		}
		return
	}

	var vd views.Data
	vd.Yield = post
	p.ShowView.Render(res, vd)
}

// Create : POST /posts
func (p *Posts) Create(res http.ResponseWriter, req *http.Request) {
	var vd views.Data
	var form PostForm
	if err := parseForm(req, &form); err != nil {
		vd.SetAlert(err)
		p.New.Render(res, vd)
		return
	}
	user := context.User(req.Context())
	post := models.Post{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := p.ps.Create(&post); err != nil {
		vd.SetAlert(err)
		p.New.Render(res, vd)
		return
	}
	fmt.Fprintln(res, post)
}

// Update : PUT /posts/:id
func (p *Posts) Update(res http.ResponseWriter, req *http.Request) {
	// TODO implement
}

// Delete : DELETE /posts/:id
func (p *Posts) Delete(res http.ResponseWriter, req *http.Request) {
	// TODO implement
}
