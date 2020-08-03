package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"nathanielwheeler.com/context"
	"nathanielwheeler.com/models"
	"nathanielwheeler.com/views"

	"github.com/gorilla/mux"
)

// ShowPost is a named route that will handle showing posts
const (
	IndexPosts = "index_posts"
	ShowPost   = "show_post"
	EditPost   = "edit_post"

	maxMultipartMem = 1 << 20 // 1 megabyte
)

// Posts will hold information about views and services
type Posts struct {
	New       *views.View
	ShowView  *views.View
	IndexView *views.View
	EditView  *views.View
	ps        models.PostsService
	r         *mux.Router
}

// NewPosts is a constructor for Posts struct
func NewPosts(ps models.PostsService, r *mux.Router) *Posts {
	return &Posts{
		New:       views.NewView("app", "posts/new"),
		ShowView:  views.NewView("app", "posts/show"),
		IndexView: views.NewView("app", "posts/index"),
		EditView:  views.NewView("app", "posts/edit"),
		ps:        ps,
		r:         r,
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
	p.ShowView.Render(res, req, vd)
}

// Index : GET /posts
func (p *Posts) Index(res http.ResponseWriter, req *http.Request) {
	posts, err := p.ps.GetAll()
	if err != nil {
		http.Error(res, "Something bad happened.", http.StatusInternalServerError)
		return
	}
	var vd views.Data
	vd.Yield = posts
	p.IndexView.Render(res, req, vd)
}

// Create : POST /posts
func (p *Posts) Create(res http.ResponseWriter, req *http.Request) {
	var vd views.Data
	var form PostForm
	if err := parseForm(req, &form); err != nil {
		vd.SetAlert(err)
		p.New.Render(res, req, vd)
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
		p.New.Render(res, req, vd)
		return
	}
	urlYear := strconv.Itoa(post.Year)
	url, err := p.r.Get(EditPost).URL("year", urlYear, "title", post.URLTitle)
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
	p.EditView.Render(res, req, vd)
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
		p.EditView.Render(res, req, vd)
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
	p.EditView.Render(res, req, vd)
}

// Upload : POST /posts/:year/:urltitle/upload
/*	- This does NOT update the path of the post. */
func (p *Posts) Upload(res http.ResponseWriter, req *http.Request) {
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

	// Parse a multipart form
	var vd views.Data
	vd.Yield = post
	err = req.ParseMultipartForm(maxMultipartMem)
	if err != nil {
		vd.SetAlert(err)
		p.EditView.Render(res, req, vd)
		return
	}

	// Create the directory for storing images
	postPath := fmt.Sprintf("images/posts/%v/%v/", post.Year, post.URLTitle)
	err = os.MkdirAll(postPath, 0755)
	if err != nil {
		vd.SetAlert(err)
		p.EditView.Render(res, req, vd)
		return
	}

	files := req.MultipartForm.File["images"]
	for _, f := range files {
		// Open uploaded files
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			p.EditView.Render(res, req, vd)
			return
		}
		defer file.Close()

		// Create a destination file
		dst, err := os.Create(postPath + f.Filename)
		if err != nil {
			vd.SetAlert(err)
			p.EditView.Render(res, req, vd)
			return
		}
		defer dst.Close()

		// Copy uploaded file data to destination
		_, err = io.Copy(dst, file)
		if err != nil {
			vd.SetAlert(err)
			p.EditView.Render(res, req, vd)
			return
		}
	}

	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Images uploaded successfully!",
	}
	p.EditView.Render(res, req, vd)
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
		p.EditView.Render(res, req, vd)
		return
	}
	url, err := p.r.Get(IndexPosts).URL()
	if err != nil {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}
	http.Redirect(res, req, url.Path, http.StatusFound)
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
