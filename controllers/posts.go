package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"nathanielwheeler.com/context"
	"nathanielwheeler.com/models"
	"nathanielwheeler.com/views"

	"github.com/gorilla/mux"
)

// Named routes.
const (
	BlogIndexRoute = "blog_index"
	BlogPostRoute  = "blog_post"
	EditPost       = "edit_post"
)

const (
	maxMultipartMem = 1 << 20 // 1 megabyte
)

// Posts will hold information about views and services
type Posts struct {
	HomeView      *views.View
	BlogPostView  *views.View
	BlogIndexView *views.View
	New           *views.View
	EditView      *views.View
	FeedView      *views.View
	ps            models.PostsService
	is            models.ImagesService
	r             *mux.Router
}

// NewPosts is a constructor for Posts struct
func NewPosts(ps models.PostsService, is models.ImagesService, r *mux.Router) *Posts {
	return &Posts{
		HomeView:      views.NewView("app", "posts/home", "posts/blog/card"),
		BlogPostView:  views.NewView("app", "posts/blog/post", "posts/blog/card"),
		BlogIndexView: views.NewView("app", "posts/blog/index"),
		New:           views.NewView("app", "posts/new"),
		EditView:      views.NewView("app", "posts/edit"),
		ps:            ps,
		is:            is,
		r:             r,
	}
}

// Home : GET /
// Needs to render the latest post
func (p *Posts) Home(res http.ResponseWriter, req *http.Request) {
	posts, err := p.ps.GetAll()
	if err != nil {
		log.Println(err)
		return
	}

	for i, post := range posts {
		err := p.ps.ParseMD(&post)
		if err != nil {
			log.Println(err)
		}
		posts[i] = post
  }

	var vd views.Data
	vd.Yield = posts
	p.HomeView.Render(res, req, vd)
}

// BlogPost : GET /blog/:filepath
func (p *Posts) BlogPost(res http.ResponseWriter, req *http.Request) {
	post, err := p.postByURL(res, req)
	if err != nil {
		// postByYearAndTitle already renders error
		return
	}
	err = p.ps.ParseMD(post)
	if err != nil {
		log.Println(err)
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

  err = p.ps.MakePostsFeed()
  if err != nil {
    vd.SetAlert(err)
    vd.RedirectAlert(res, req, "/home", http.StatusFound, *vd.Alert)
    return
  }
	p.BlogIndexView.Render(res, req, vd)
}

// PostForm will hold information for creating a new post
type PostForm struct {
	Title    string `schema:"title"`
	URLPath  string `schema:"urlpath"`
	FilePath string `schema:"filepath"`
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
	if user.IsAdmin != true {
		http.Error(res, "You do not have permission to create a post", http.StatusForbidden)
		return
	}
	post := models.Post{
		Title:    form.Title,
		URLPath:  form.URLPath,
		FilePath: "public/markdown/" + form.FilePath + ".md",
	}
	if err := p.ps.Create(&post); err != nil {
		vd.SetAlert(err)
		p.New.Render(res, req, vd)
		return
  }
	url, err := p.r.Get(BlogPostRoute).URL("id", fmt.Sprintf("%v", post.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(res, req, "/blog", http.StatusFound)
		return
  }
  err = p.ps.MakePostsFeed()
  if err != nil {
    vd.SetAlert(err)
    vd.RedirectAlert(res, req, url.Path, http.StatusFound, *vd.Alert)
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
	if user.IsAdmin != true {
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
	if user.IsAdmin != true {
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
	if user.IsAdmin != true {
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
	if user.IsAdmin != true {
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
	if user.IsAdmin != true {
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
	url, err := p.r.Get(EditPost).URL("id", fmt.Sprintf("%v", post.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(res, req, "/posts/index", http.StatusFound)
		return
	}
	http.Redirect(res, req, url.Path, http.StatusFound)
}

// #region HELPERS

func (p *Posts) postByURL(res http.ResponseWriter, req *http.Request) (*models.Post, error) {
	vars := mux.Vars(req)
	urlpath := vars["urlpath"]
	post, err := p.ps.ByURL(urlpath)
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
	return post, nil
}

func (p *Posts) postByLatest(res http.ResponseWriter, req *http.Request) (*models.Post, error) {
	post, err := p.ps.ByLatest()
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(res, "Post not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(res, "Well, that wasn't supposed to happen", http.StatusInternalServerError)
		}
		return nil, err
	}
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
	return post, nil
}

// #endregion
