package main

import (
	"fmt"
	"net/http"

	"nathanielwheeler.com/controllers"
	"nathanielwheeler.com/middleware"
	"nathanielwheeler.com/models"
	"nathanielwheeler.com/rand"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

/* TODO
- Decouple User middleware from routes (like with RequireUser)
*/

func main() {
	cfg := LoadConfig()
	dbCfg := cfg.Database

	// Initialize services
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionString()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithPosts(),
		models.WithImages(),
	)
	defer services.Close()
	services.AutoMigrate()

	// Router Initialization
	r := mux.NewRouter()

	// Initialize controllers
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	postsC := controllers.NewPosts(services.Posts, services.Images, r)

	// Middleware
	userMw := middleware.User{UserService: services.User}
	requireUserMw := middleware.RequireUser{}

	// CSRF Protection
	b, err := rand.Bytes(cfg.CSRFBytes)
	if err != nil {
		panic(err)
	}
	csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProd()))

	// Public Routes
	publicHandler := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/images/").
		Handler(publicHandler)
	r.PathPrefix("/stylesheets/").
		Handler(publicHandler)
	r.PathPrefix("/markdown/").
		Handler(publicHandler)

	// Statics Routes
	r.Handle("/resume",
		staticC.Resume).
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
    requireUserMw.ApplyFn(usersC.Logout)).
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
    Name(controllers.BlogIndexRoute)
  r.HandleFunc(`/blog/{urlpath:[a-zA-Z0-9\/\-_~.]+}`,
    postsC.BlogPost).
    Methods("GET").
    Name(controllers.BlogPostRoute)
  //    API / Admin
	r.HandleFunc("/posts",
		requireUserMw.ApplyFn(postsC.Create)).
		Methods("POST")
	r.Handle("/posts/new",
		requireUserMw.Apply(postsC.New)).
		Methods("GET")
	r.HandleFunc("/posts/{id:[0-9]+}/edit",
		requireUserMw.ApplyFn(postsC.Edit)).
		Methods("GET").
		Name(controllers.EditPost)
	r.HandleFunc("/posts/{id:[0-9]+}/update",
		requireUserMw.ApplyFn(postsC.Update)).
		Methods("POST")
	r.HandleFunc("/posts/{id:[0-9]+}/delete",
		requireUserMw.ApplyFn(postsC.Delete)).
		Methods("POST")
		//    Images
	r.HandleFunc("/posts/{id:[0-9]+}/upload",
		requireUserMw.ApplyFn(postsC.ImageUpload)).
		Methods("POST")
	r.HandleFunc("/posts/{id:[0-9]+}/image/{filename}/delete",
		requireUserMw.ApplyFn(postsC.ImageDelete)).
		Methods("POST")

	// Start that server!
	port := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("Now listening on %s...\n", port)
	http.ListenAndServe(port, csrfMw(userMw.Apply(r)))
}
