package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"nathanielwheeler.com/context"
	"nathanielwheeler.com/models"
	"nathanielwheeler.com/views"

	"github.com/gorilla/mux"
)

// Named routes.
const (
	BlogIndexRoute = "blog_index"
	BlogPostRoute  = "blog_post"
	EditPost  = "edit_post"
)

const (
	maxMultipartMem = 1 << 20 // 1 megabyte
)

// Posts will hold information about views and services
type Posts struct {
	New           *views.View
	BlogPostView  *views.View
	BlogIndexView *views.View
	EditView      *views.View
	ps            models.PostsService
	is            models.ImagesService
	r             *mux.Router
}

// NewPosts is a constructor for Posts struct
func NewPosts(ps models.PostsService, is models.ImagesService, r *mux.Router) *Posts {
	return &Posts{
		New:           views.NewView("app", "posts/new"),
		EditView:      views.NewView("app", "posts/edit"),
		BlogPostView:  views.NewView("app", "posts/blog/post"),
		BlogIndexView: views.NewView("app", "posts/blog/index"),
		ps:            ps,
		is:            is,
		r:             r,
	}
}

// PostForm will hold information for creating a new post
type PostForm struct {
	Title string `schema:"title"`
}

// BlogPost : GET /blog/:year/:urltitle
func (p *Posts) BlogPost(res http.ResponseWriter, req *http.Request) {
	post, err := p.postByYearAndTitle(res, req)
	if err != nil {
		// postByYearAndTitle already renders error
		return
	}
	var vd views.Data
	vd.Yield = post
	p.BlogPostView.Render(res, req, vd)
}

// BlogIndex : GET /blog
func (p *Posts) BlogIndex(res http.ResponseWriter, req *http.Request) {
	posts, err := p.ps.GetAll()
	if err != nil {
		log.Println(err)
		http.Error(res, "Something bad happened.", http.StatusInternalServerError)
		return
	}
	var vd views.Data
	vd.Yield = posts
	p.BlogIndexView.Render(res, req, vd)
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
		log.Println(err)
		http.Redirect(res, req, "/posts/index", http.StatusFound)
		return
	}
	http.Redirect(res, req, url.Path, http.StatusFound)
}

// Edit : POST /posts/:id/update
func (p *Posts) Edit(res http.ResponseWriter, req *http.Request) {
	post, err := p.postByID(res, req)
	if err != nil {
		// error handled by postByID
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

// Update : POST /posts/:id/update
/*  - This does NOT update the path of the post. */
func (p *Posts) Update(res http.ResponseWriter, req *http.Request) {
	post, err := p.postByID(res, req)
	if err != nil {
		// implemented by postByID
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

// Delete : POST /posts/:id/delete
// TODO cascade delete all images within post
func (p *Posts) Delete(res http.ResponseWriter, req *http.Request) {
	post, err := p.postByID(res, req)
	if err != nil {
		// postByID renders error
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
	url, err := p.r.Get(BlogIndexRoute).URL()
	if err != nil {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}
	http.Redirect(res, req, url.Path, http.StatusFound)
}

// ImageUpload : POST /posts/:id/upload
func (p *Posts) ImageUpload(res http.ResponseWriter, req *http.Request) {
	post, err := p.postByID(res, req)
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

		// Create image
		err = p.is.Create(post.ID, file, f.Filename)
		if err != nil {
			vd.SetAlert(err)
			p.EditView.Render(res, req, vd)
			return
		}
	}

	url, err := p.r.Get(EditPost).URL("id", fmt.Sprintf("%v", post.ID))
	if err != nil {
		http.Redirect(res, req, "/posts", http.StatusFound)
		return
	}
	http.Redirect(res, req, url.Path, http.StatusFound)
}

// ImageDelete /posts/:id/images/:filename/delete
func (p *Posts) ImageDelete(res http.ResponseWriter, req *http.Request) {
	post, err := p.postByID(res, req)
	if err != nil {
		return
	}
	user := context.User(req.Context())
	if post.UserID != user.ID {
		http.Error(res, "You do not have permission to edit this post or image", http.StatusForbidden)
		return
	}
	filename := mux.Vars(req)["filename"]
	i := models.Image{
		Filename: filename,
		PostID:   post.ID,
	}
	err = p.is.Delete(&i)
	if err != nil {
		var vd views.Data
		vd.Yield = post
		vd.SetAlert(err)
		p.EditView.Render(res, req, vd)
		return
	}
	url, err := p.r.Get(EditPost).URL("year", fmt.Sprintf("%v", post.Year), "title", fmt.Sprintf("%v", post.URLTitle))
	if err != nil {
		log.Println(err)
		http.Redirect(res, req, "/posts/index", http.StatusFound)
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
		log.Println(err)
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
			log.Println(err)
			http.Error(res, "Something bad happened.", http.StatusInternalServerError)
		}
		return nil, err
	}
	post = p.getImages(post)
	return post, nil
}

func (p *Posts) postByID(res http.ResponseWriter, req *http.Request) (*models.Post, error) {
	idVar := mux.Vars(req)["id"]
	id, err := strconv.Atoi(idVar)
	if err != nil {
		log.Println(err)
		http.Error(res, "Invalid post ID", http.StatusNotFound)
		return nil, err
	}
	post, err := p.ps.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(res, "Post not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(res, "Something bad happened.", http.StatusInternalServerError)
		}
		return nil, err
	}
	post = p.getImages(post)
	return post, nil
}

// Get images from ImageService and attach to post.
func (p *Posts) getImages(post *models.Post) *models.Post {
	images, err := p.is.ByPostID(post.ID)
	if err != nil {
		log.Println(err)
	}
	post.Images = images
	return post
}

// #endregion
