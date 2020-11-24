package server

import (
	"net/http"

	"nathanielwheeler.com/server/handlers"
	"nathanielwheeler.com/server/middleware"
)

func (s *server) routes() {
	r := s.router

	// Middleware
	// userMw := middleware.User{UserService: s.services.User}
	// b, err := util.Bytes(s.config.CSRFBytes)
	// if err != nil {
	// 	panic(err)
	// }
	// csrfMw := csrf.Protect(b, csrf.Secure(s.isProd()))

	// Initialize handlers
	staticH := handlers.NewStatic()
	usersH := handlers.NewUsers(s.services.User)
	postsH := handlers.NewPosts(s.services.Posts, s.services.Images, r)

	// Fileserver Routes
	publicHandler := http.FileServer(http.Dir("./client/public/"))
	r.PathPrefix("/images/").
		Handler(publicHandler)
	r.PathPrefix("/assets/").
		Handler(publicHandler)
	r.PathPrefix("/markdown/").
		Handler(publicHandler)
	r.PathPrefix("/feeds/").
		Handler(publicHandler)

	// Statics Routes
	r.HandleFunc("/resume",
		s.handleTemplate(nil, "app", "static/resume")).
		Methods("GET")
	r.HandleFunc("/prototypes/theme-system",
		s.handleTemplate(nil, "app", "static/prototypes/theme-system")).
		Methods("GET")

	// User Routes
	r.HandleFunc("/register",
		s.handleTemplate(nil, "app", "users/register")).
		Methods("GET")
	r.HandleFunc("/register",
		usersH.handleRegister).
		Methods("POST")
	r.Handle("/login",
		usersH.LoginView).
		Methods("GET")
	r.HandleFunc("/login",
		usersH.Login).
		Methods("POST")
	r.Handle("/logout",
		middleware.ApplyFn(usersH.Logout)).
		Methods("POST")
	r.HandleFunc("/cookietest",
		usersH.CookieTest).
		Methods("GET")

	// Post Routes
	//    Blog
	r.HandleFunc("/",
		postsH.Home).
		Methods("GET")
	r.HandleFunc("/blog",
		postsH.BlogIndex).
		Methods("GET").
		Name(handlers.BlogIndexRoute)
	r.HandleFunc(`/blog/{urlpath:[a-zA-Z0-9\/\-_~.]+}`,
		postsH.BlogPost).
		Methods("GET").
		Name(handlers.BlogPostRoute)
		//    API / Admin
	r.HandleFunc("/posts",
		middleware.ApplyFn(postsH.Create)).
		Methods("POST")
	r.Handle("/posts/new",
		middleware.Apply(postsH.New)).
		Methods("GET")
	r.HandleFunc("/posts/{id:[0-9]+}/edit",
		middleware.ApplyFn(postsH.Edit)).
		Methods("GET").
		Name(handlers.EditPost)
	r.HandleFunc("/posts/{id:[0-9]+}/update",
		middleware.ApplyFn(postsH.Update)).
		Methods("POST")
	r.HandleFunc("/posts/{id:[0-9]+}/delete",
		middleware.ApplyFn(postsH.Delete)).
		Methods("POST")
		//    Images
	r.HandleFunc("/posts/{id:[0-9]+}/upload",
		middleware.ApplyFn(postsH.ImageUpload)).
		Methods("POST")
	r.HandleFunc("/posts/{id:[0-9]+}/image/{filename}/delete",
		middleware.ApplyFn(postsH.ImageDelete)).
		Methods("POST")
}
