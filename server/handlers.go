package server

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"nathanielwheeler.com/server/services"
	"nathanielwheeler.com/server/util"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func (s *server) viewData(data interface{}, layout string, files ...string) Data {
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}
	vd.View = s.newView(layout, files...)
	return vd
}

func (s *server) handleTemplate(data interface{}, layout string, files ...string) http.HandlerFunc {
	var (
		vd   Data
		init sync.Once
	)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		init.Do(func() {
			s.viewData(nil, layout, files...)
		})

		// Check cookie for alerts
		if alert := getAlert(r); alert != nil {
			vd.Alert = alert
			clearAlert(w)
		}

		// Lookup and set the user to the User field
		vd.User = util.User(r.Context())
		var buf bytes.Buffer

		// Create CSRF field using current http request and add it onto the template FuncMap.
		csrfField := csrf.TemplateField(r)

		tpl := vd.View.Template.Funcs(template.FuncMap{
			"csrfField": func() template.HTML {
				return csrfField
			},
		})

		err := tpl.ExecuteTemplate(&buf, vd.View.Layout, vd)
		if err != nil {
			http.Error(w, `Something went wrong, please try again.  If the problem persists, please contact me directly at "nathan@mailftp.com"`, http.StatusInternalServerError)
			return
		}
		io.Copy(w, &buf)
	}
}

// USERS

// RegistrationForm is used to transform a webform into a registration request
type RegistrationForm struct {
	Email    string `schema:"email"`
	Name     string `schema:"name"`
	Password string `schema:"password"`
}

// handleRegister : POST /register
// — Used to process the signup form when a user tries to create a new user account
func (s *server) handleRegister() http.HandlerFunc {

	var vd Data
	var form RegistrationForm
	return func(w http.ResponseWriter, r *http.Request) {
		if err := parseForm(r, &form); err != nil {
			vd.SetAlert(err)
			u.RegisterView.Render(w, r, vd)
			return
		}
		user := services.User{
			Name:     form.Name,
			Email:    form.Email,
			Password: form.Password,
		}
		if err := u.us.Create(&user); err != nil {
			vd.SetAlert(err)
			u.RegisterView.Render(w, r, vd)
			return
		}
		err := s.signIn(w, &user)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/cookietest", http.StatusFound)
	}
}

// LoginForm is used to transform a webform into a login request
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Login : POST /login
// — Used to process the login form when a user tries to log in as an existing user
func (s *server) Login(res http.ResponseWriter, req *http.Request) {
	var vd views.Data
	var form LoginForm
	if err := parseForm(req, &form); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(res, req, vd)
		return
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case services.ErrNotFound:
		case services.ErrPasswordInvalid:
			vd.AlertError("Invalid email and/or password.")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(res, req, vd)
		return
	}

	err = u.signIn(res, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(res, req, vd)
		return
	}
	http.Redirect(res, req, "/cookietest", http.StatusFound)
}

// Logout : POST /logout
// — Used to process the logout form when a user chooses to logout.
func (s *server) Logout(res http.ResponseWriter, req *http.Request) {
	// Expire user's cookie
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(res, &cookie)
	// Update user with a new remember token
	user := util.User(req.Context())
	token, _ := util.RememberToken() // Ignoring errors because because unlikely and not much to do about it.
	user.Remember = token
	u.us.Update(user)

	http.Redirect(res, req, "/", http.StatusFound)
}

// signIn is used to sign the given user in via cookies
func (s *server) signIn(res http.ResponseWriter, user *services.User) error {
	if user.Remember == "" {
		token, err := util.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(res, &cookie)
	return nil
}

// CookieTest is used to display cookies set on the current user
func (s *server) CookieTest(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("remember_token")
	if err != nil {
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}
	_, err = u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}
	http.Redirect(res, req, "/", http.StatusFound)
}

// POSTS

// Named routes.
const (
	BlogIndexRoute = "blog_index"
	BlogPostRoute  = "blog_post"
	EditPost       = "edit_post"
)

const (
	maxMultipartMem = 1 << 20 // 1 megabyte
)

// PostViews will hold information about views and services
type PostViews struct {
	HomeView      *View
	BlogPostView  *View
	BlogIndexView *View
	New           *View
	EditView      *View
	FeedView      *View
}

func (s *server) newPostViews() *PostViews {
	return &PostViews{
		HomeView:      NewView("app", "posts/home", "posts/blog/card"),
		BlogPostView:  NewView("app", "posts/blog/post", "posts/blog/card"),
		BlogIndexView: NewView("app", "posts/blog/index"),
		New:           NewView("app", "posts/new"),
		EditView:      NewView("app", "posts/edit"),
	}
}

// Home : GET /
// Needs to render the latest post
func (s *server) Home(res http.ResponseWriter, req *http.Request) {
	posts, err := s.services.Posts.GetAll()
	if err != nil {
		s.logErr("home error")
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
func (s *server) BlogPost(res http.ResponseWriter, req *http.Request) {
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
func (s *server) BlogIndex(res http.ResponseWriter, req *http.Request) {
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
func (s *server) Create(res http.ResponseWriter, req *http.Request) {
	var vd views.Data
	var form PostForm
	if err := parseForm(req, &form); err != nil {
		vd.SetAlert(err)
		p.New.Render(res, req, vd)
		return
	}
	user := util.User(req.Context())
	if user.IsAdmin != true {
		http.Error(res, "You do not have permission to create a post", http.StatusForbidden)
		return
	}
	post := services.Post{
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
	// Allows me to test pages locally without updating feeds
	if p.ps.IsProduction() {
		err = p.ps.MakePostsFeed()
		if err != nil {
			vd.SetAlert(err)
			vd.RedirectAlert(res, req, url.Path, http.StatusFound, *vd.Alert)
			return
		}
	}
	http.Redirect(res, req, url.Path, http.StatusFound)
}

// Edit : POST /posts/:id/update
func (s *server) Edit(res http.ResponseWriter, req *http.Request) {
	post, err := p.postByID(res, req)
	if err != nil {
		// error handled by postByID
		return
	}
	user := util.User(req.Context())
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
func (s *server) Update(res http.ResponseWriter, req *http.Request) {
	post, err := p.postByID(res, req)
	if err != nil {
		// implemented by postByID
		return
	}
	user := util.User(req.Context())
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
func (s *server) Delete(res http.ResponseWriter, req *http.Request) {
	post, err := p.postByID(res, req)
	if err != nil {
		// postByID renders error
		return
	}
	user := util.User(req.Context())
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

func (s *server) postByURL(res http.ResponseWriter, req *http.Request) (*services.Post, error) {
	vars := mux.Vars(req)
	urlpath := vars["urlpath"]
	post, err := p.ps.ByURL(urlpath)
	if err != nil {
		switch err {
		case services.ErrNotFound:
			http.Error(res, "Post not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(res, "Something bad happened.", http.StatusInternalServerError)
		}
		return nil, err
	}
	return post, nil
}

func (s *server) postByLatest(res http.ResponseWriter, req *http.Request) (*services.Post, error) {
	post, err := p.ps.ByLatest()
	if err != nil {
		switch err {
		case services.ErrNotFound:
			http.Error(res, "Post not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(res, "Well, that wasn't supposed to happen", http.StatusInternalServerError)
		}
		return nil, err
	}
	return post, nil
}

func (s *server) postByID(res http.ResponseWriter, req *http.Request) (*services.Post, error) {
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
		case services.ErrNotFound:
			http.Error(res, "Post not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(res, "Something bad happened.", http.StatusInternalServerError)
		}
		return nil, err
	}
	return post, nil
}

// HELPERS

func parseForm(req *http.Request, dest interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	decoder := schema.NewDecoder()
	// IgnoreUnknownKeys so that we can use CSRF protection in our forms.
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(dest, req.PostForm); err != nil {
		return err
	}

	return nil
}
