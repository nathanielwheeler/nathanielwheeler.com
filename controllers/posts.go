package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	EditView *views.View
	ps       models.PostsService
	r        *mux.Router
}

// NewPosts is a constructor for Posts struct
func NewPosts(ps models.PostsService, r *mux.Router) *Posts {
	return &Posts{
		New:      views.NewView("app", "posts/new"),
		ShowView: views.NewView("app", "posts/show"),
		EditView: views.NewView("app", "posts/edit"),
		ps:       ps,
		r:        r,
	}
}

// PostForm will hold information for creating a new post
type PostForm struct {
	Title string `schema:"title"`
}

// Show : GET /posts/:year/:urltitle
func (p *Posts) Show(res http.ResponseWriter, req *http.Request) {
	post, err := p.postByYearAndTitle(res, req)
	if err != nil {
		// postByYearAndTitle already renders error
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
		Title:    form.Title,
		URLTitle: strings.Replace(form.Title, " ", "_", -1),
		UserID:   user.ID,
		Year:     time.Now().Year(),
	}
	if err := p.ps.Create(&post); err != nil {
		vd.SetAlert(err)
		p.New.Render(res, vd)
		return
	}
	urlYear := strconv.Itoa(post.Year)
	url, err := p.r.Get(ShowPost).URL("year", urlYear, "title", post.URLTitle)
	if err != nil {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}
	http.Redirect(res, req, url.Path, http.StatusFound)
}

// Edit : POST /posts/:year/:urltitle/update
func (p *Posts) Edit(res http.ResponseWriter, req *http.Request) {
	post, err := p.postByYearAndTitle(res, req)
	if err != nil {
		// error handled by postByYearAndTitle
		return
	}
	user := context.User(req.Context())
	if post.UserID != user.ID {
		http.Error(res, "You do not have permission to edit this post", http.StatusForbidden)
		return
	}
	var vd views.Data
	vd.Yield = post
	p.EditView.Render(res, vd)
}

// Update : POST /posts/:year/:urltitle/update
/*	- This does NOT update the path of the post. */
func (p *Posts) Update(res http.ResponseWriter, req *http.Request) {
	post, err := p.postByYearAndTitle(res, req)
	if err != nil {
		// implemented by postByYearAndTitle
		return
	}
	user := context.User(req.Context())
	if post.UserID != user.ID {
		http.Error(res, "You do not have permission to edit this post", http.StatusForbidden)
		return
	}

	var vd views.Data
	vd.Yield = post
	var form PostForm
	if err := parseForm(req, &form); err != nil {
		vd.SetAlert(err)
		p.EditView.Render(res, vd)
		return
	}
	post.Title = form.Title
	err = p.ps.Update(post)
	if err != nil {
		vd.SetAlert(err)
	} else {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Post updated successfully!",
		}
	}
	p.EditView.Render(res, vd)
}

// Delete : POST /posts/:year/:urltitle/delete
func (p *Posts) Delete(res http.ResponseWriter, req *http.Request) {
	post, err := p.postByYearAndTitle(res, req)
	if err != nil {
		// postByYearAndTitle renders error
		return
	}
	user := context.User(req.Context())
	if post.UserID != user.ID {
		http.Error(res, "You do not have permission to edit this post", http.StatusForbidden)
		return
	}
	var vd views.Data
	err = p.ps.Delete(post.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = post
		p.EditView.Render(res, vd)
		return
	}
	fmt.Fprintln(res, "Successfully deleted!")
}

// #region HELPERS

func (p *Posts) postByYearAndTitle(res http.ResponseWriter, req *http.Request) (*models.Post, error) {
	vars := mux.Vars(req)
	yearVar := vars["year"]
	year, err := strconv.Atoi(yearVar)
	if err != nil {
		http.Error(res, "Invalid post URL", http.StatusNotFound)
		return nil, err
	}
	urlTitle := vars["title"]
	post, err := p.ps.ByYearAndTitle(year, urlTitle)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(res, "Post not found", http.StatusNotFound)
		default:
			http.Error(res, "Something bad happened.", http.StatusInternalServerError)
		}
		return nil, err
	}
	return post, nil
}

// #endregion
