package server

import (
	"net/http"

	"github.com/gorilla/csrf"
	"nathanielwheeler.com/server/handlers"
	"nathanielwheeler.com/server/middleware"
	"nathanielwheeler.com/server/services"
	"nathanielwheeler.com/server/util"
)

func (s *server) routes() {
	r := s.router

	// Middleware
	userMw := middleware.User{UserService: services.User}
	b, err := util.Bytes(s.config.CSRFBytes)
	if err != nil {
		panic(err)
	}
	csrfMw := csrf.Protect(b, csrf.Secure(s.isProd()))

	// Initialize handlers
	staticC := handlers.NewStatic()
	usersC := handlers.NewUsers(services.User)
	postsC := handlers.NewPosts(services.Posts, services.Images, r)

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
	r.Handle("/resume",
		staticC.Resume).
		Methods("GET")
	r.Handle("/prototypes/theme-system",
		staticC.PrototypeThemeSystem).
		Methods("GET")

	// User Routes
	r.HandleFunc("/register",
		usersC.Registration).
		Methods("GET")
	r.HandleFunc("/register",
		usersC.Register).
		Methods("POST")
	r.Handle("/login",
		usersC.LoginView).
		Methods("GET")
	r.HandleFunc("/login",
		usersC.Login).
		Methods("POST")
	r.Handle("/logout",
		middleware.ApplyFn(usersC.Logout)).
		Methods("POST")
	r.HandleFunc("/cookietest",
		usersC.CookieTest).
		Methods("GET")

	// Post Routes
	//    Blog
	r.HandleFunc("/",
		postsC.Home).
		Methods("GET")
	r.HandleFunc("/blog",
		postsC.BlogIndex).
		Methods("GET").
		Name(handlers.BlogIndexRoute)
	r.HandleFunc(`/blog/{urlpath:[a-zA-Z0-9\/\-_~.]+}`,
		postsC.BlogPost).
		Methods("GET").
		Name(handlers.BlogPostRoute)
		//    API / Admin
	r.HandleFunc("/posts",
		middleware.ApplyFn(postsC.Create)).
		Methods("POST")
	r.Handle("/posts/new",
		middleware.Apply(postsC.New)).
		Methods("GET")
	r.HandleFunc("/posts/{id:[0-9]+}/edit",
		middleware.ApplyFn(postsC.Edit)).
		Methods("GET").
		Name(handlers.EditPost)
	r.HandleFunc("/posts/{id:[0-9]+}/update",
		middleware.ApplyFn(postsC.Update)).
		Methods("POST")
	r.HandleFunc("/posts/{id:[0-9]+}/delete",
		middleware.ApplyFn(postsC.Delete)).
		Methods("POST")
		//    Images
	r.HandleFunc("/posts/{id:[0-9]+}/upload",
		middleware.ApplyFn(postsC.ImageUpload)).
		Methods("POST")
	r.HandleFunc("/posts/{id:[0-9]+}/image/{filename}/delete",
		middleware.ApplyFn(postsC.ImageDelete)).
		Methods("POST")
}
