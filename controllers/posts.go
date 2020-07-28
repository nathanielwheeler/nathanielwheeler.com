package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"nathanielwheeler.com/context"
	"nathanielwheeler.com/models"
	"nathanielwheeler.com/views"

	"github.com/gorilla/mux"
)

// ShowPost is a named route that will handle showing posts
const ShowPost = "show_post"

// Posts will hold information about views and services
type Posts struct {
	New      *views.View
	ShowView *views.View
	ps       models.PostsService
	r        *mux.Router
}

// NewPosts is a constructor for Posts struct
func NewPosts(ps models.PostsService, r *mux.Router) *Posts {
	return &Posts{
		New:      views.NewView("app", "posts/new"),
		ShowView: views.NewView("app", "posts/show"),
		ps:       ps,
		r:        r,
	}
}

// PostForm will hold information for creating a new post
type PostForm struct {
	Title string `schema:"title"`
}

// Show : GET /posts/:year/:title
func (p *Posts) Show(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	// TODO use ByYearAndTitle instead of ByID
	yearVar := vars["year"]
	year, err := strconv.Atoi(yearVar)
	if err != nil {
		http.Error(res, "Invalid post URL", http.StatusNotFound)
		return
	}
	titleVar := vars["title"]
	title := strings.Replace(titleVar, "_", " ", -1)
	post, err := p.ps.ByYearAndTitle(year, title)

	// Alternate way to show post by ID instead of Year and Title
	/* 	idStr := vars["id"]
	   	id, err := strconv.Atoi(idStr)
	   	if err != nil {
	   		http.Error(res, "Invalid post ID", http.StatusNotFound)
	   		return
	   	}
	   	post, err := p.ps.ByID(uint(id)) */

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

	// Redirect to new post
	// TODO implement ByYearAndTitle
	urlTitle := strings.Replace(post.Title, " ", "_", -1)
	urlYear := strconv.Itoa(post.CreatedAt.Year())
	url, err := p.r.Get(ShowPost).URL("year", urlYear, "title", urlTitle)
	// If I want to use ID instead...
	// url, err := p.r.Get(ShowPost).URL("id", strconv.Itoa(int(post.ID)))
	if err != nil {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}
	http.Redirect(res, req, url.Path, http.StatusFound)
}
