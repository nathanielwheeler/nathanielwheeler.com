package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"nathanielwheeler.com/controllers"
	"nathanielwheeler.com/middleware"
	"nathanielwheeler.com/models"
	"nathanielwheeler.com/rand"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const (
	port   = ":3000"
	isProd = false
)

/* TODO
- Decouple User middleware from routes (like with RequireUser)
*/

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

func main() {
	// Start up database connection
	dbEnv := getDBEnv()
	psqlConnectionStr := fmt.Sprintf(
		"host=%s port=%s user=%s password='%s' dbname=%s sslmode=disable",
		dbEnv.host, dbEnv.port, dbEnv.user, dbEnv.password, dbEnv.name,
	)

	// Initialize services
	services, err := models.NewServices(psqlConnectionStr)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	// Router Initilization
	r := mux.NewRouter()

	// Initialize controllers
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	postsC := controllers.NewPosts(services.Posts, services.Images, r)

	// Middleware
	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{}
	// CSRF Protection
	b, err := rand.Bytes(32)
	if err != nil {
		panic(err)
	}
	csrfMw := csrf.Protect(b, csrf.Secure(isProd))

	// Public Routes
	publicHandler := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/images/").
		Handler(publicHandler)
	r.PathPrefix("/assets/").
		Handler(publicHandler)

	// Statics Routes
	r.Handle("/",
		staticC.Home).
		Methods("GET")
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
	r.HandleFunc("/cookietest",
		usersC.CookieTest).
		Methods("GET")

	// Post Routes
	r.HandleFunc("/posts",
		requireUserMw.ApplyFn(postsC.Create)).
		Methods("POST")
	// FIXME: I can't figure out why Index will render for GET "/posts/index" but not GET "/posts"
	r.HandleFunc("/posts/index",
		postsC.Index).
		Methods("GET").
		Name(controllers.IndexPosts)
	r.Handle("/posts/new",
		requireUserMw.Apply(postsC.New)).
		Methods("GET")
	r.HandleFunc("/posts/{year:20[0-9]{2}}/{title}",
		postsC.Show).
		Methods("GET").
		Name(controllers.ShowPost)
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
	// TODO Update references to old route
	r.HandleFunc("/posts/{id:[0-9]+}/upload",
		requireUserMw.ApplyFn(postsC.ImageUpload)).
		Methods("POST")
	r.HandleFunc("/posts/{id:[0-9]+}/image/{filename}/delete",
		requireUserMw.ApplyFn(postsC.ImageDelete)).
		Methods("POST")

	// Start that server!
	fmt.Println("Now listening on", port)
	http.ListenAndServe(port, csrfMw(userMw.Apply(r)))
}

// #region DB HELPERS

type dbEnv struct {
	host, user, password, port, name string
}

func getDBEnv() dbEnv {
	return dbEnv{
		host:     checkDBEnv("host"),
		user:     checkDBEnv("user"),
		password: checkDBEnv("password"),
		port:     checkDBEnv("port"),
		name:     checkDBEnv("name"),
	}
}

func checkDBEnv(str string) string {
	str, exists := os.LookupEnv("DB_" + strings.ToUpper(str))
	if !exists {
		panic(".env is missing environment variable: '" + str + "'")
	}
	return str
}

// #endregion
