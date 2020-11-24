package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"nathanielwheeler.com/server/services"

	"github.com/gorilla/mux"
)

// Run starts up the server
func Run() (err error) {
	s := newServer()

	// Start that server!
	port := fmt.Sprintf(":%d", s.config.Port)
	fmt.Printf("Now listening on %s...\n", port)
	http.ListenAndServe(port, csrfMw(userMw.Apply(r)))

	return err
}

type server struct {
	config   *config
	logger   *log.Logger
	router   *mux.Router
	services *services.Services
}

func newServer() *server {
	s := server{
		config: loadConfig(),
		logger: log.New(os.Stdout, "server: ", log.Lshortfile),
		router: mux.NewRouter(),
	}

	// Initialize services
	err := services.NewServices(
		services.WithGorm(s.dialect(), s.connectionString()),
		services.WithLogMode(!s.isProd()),
		services.WithUser(s.config.Pepper, s.config.HMACKey),
		services.WithPosts(s.isProd()),
		services.WithImages(),
	)
	if err != nil {
		
	}
	defer s.services.Close()
	s.services.AutoMigrate()

	s.routes()
	return &s
}
