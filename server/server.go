package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"nathanielwheeler.com/server/middleware"
	"nathanielwheeler.com/server/services"
	"nathanielwheeler.com/server/util"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

// Run starts up the server
func Run() (err error) {
	s := newServer()

	// Initialize services
	services, err := services.NewServices(
		services.WithGorm(s.dialect(), s.connectionString()),
		services.WithLogMode(!s.isProd()),
		services.WithUser(s.config.Pepper, s.config.HMACKey),
		services.WithPosts(s.isProd()),
		services.WithImages(),
	)
	defer services.Close()
	services.AutoMigrate()

	// Router Initialization
	r := mux.NewRouter()

	// Middleware
	userMw := middleware.User{UserService: services.User}

	// CSRF Protection
	b, err := util.Bytes(s.config.CSRFBytes)
	if err != nil {
		panic(err)
	}
	csrfMw := csrf.Protect(b, csrf.Secure(s.isProd()))

	// Start that server!
	port := fmt.Sprintf(":%d", s.config.Port)
	fmt.Printf("Now listening on %s...\n", port)
	http.ListenAndServe(port, csrfMw(userMw.Apply(r)))

	return err
}

type server struct {
	logger *log.Logger
	router *mux.Router
	config *config
}

func newServer() *server {
	s := server{
		logger: log.New(os.Stdout, "server: ", log.Lshortfile),
		config: loadConfig(),
	}
	s.routes()
	return &s
}
