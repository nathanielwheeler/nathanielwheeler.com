package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"nathanielwheeler.com/middleware"
	"nathanielwheeler.com/controllers"
	"nathanielwheeler.com/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type env struct {
	host, user, password, port, name string
}

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

func main() {
	// Start up database connection
	dbEnv := getEnv()
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

	// Initialize controllers
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	postsC := controllers.NewPosts(services.Posts)

	// Middleware
	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}
	newPost := requireUserMw.Apply(postsC.New)
	createPost := requireUserMw.ApplyFn(postsC.Create)

	// Route Handling
	r := mux.NewRouter()
	//		Statics
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/resume", staticC.Resume).Methods("GET")
	//		Users
	r.HandleFunc("/register", usersC.Registration).Methods("GET") // Consider making this a view to match LoginView
	r.HandleFunc("/register", usersC.Register).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	//		Posts
	r.Handle("/posts/new", newPost).Methods("GET")
	r.HandleFunc("posts", createPost).Methods("POST")

	// Start that server!
	http.ListenAndServe(":3000", r)
}

// #region DB HELPERS

func getEnv() env {
	return env{
		host:     checkEnv("host"),
		user:     checkEnv("user"),
		password: checkEnv("password"),
		port:     checkEnv("port"),
		name:     checkEnv("name"),
	}
}

func checkEnv(str string) string {
	str, exists := os.LookupEnv("DB_" + strings.ToUpper(str))
	if !exists {
		panic(".env is missing environment variable: '" + str + "'")
	}
	return str
}

// #endregion
