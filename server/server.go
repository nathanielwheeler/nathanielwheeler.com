package server

import (
	"nathanielwheeler.com/server/handlers"
	"nathanielwheeler.com/server/services"
)

// Run starts up the server
func Run() (err error) {
	cfg := config.LoadConfig()
	dbCfg := cfg.Database

	// Initialize services
	services, err := services.NewServices(
		services.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionString()),
		services.WithLogMode(!cfg.IsProd()),
		services.WithUser(cfg.Pepper, cfg.HMACKey),
		services.WithPosts(cfg.IsProd()),
		services.WithImages(),
	)
	defer services.Close()
	services.AutoMigrate()

	// Router Initialization
	r := mux.NewRouter()

	// Initialize controllers
	staticC := handlers.NewStatic()
	usersC := handlers.NewUsers(services.User)
	postsC := handlers.NewPosts(services.Posts, services.Images, r)

	// Middleware
	userMw := middleware.User{UserService: services.User}
	requireUserMw := middleware.RequireUser{}

	// CSRF Protection
	b, err := rand.Bytes(cfg.CSRFBytes)
	if err != nil {
		panic(err)
	}
	csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProd()))

	// Start that server!
	port := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("Now listening on %s...\n", port)
	http.ListenAndServe(port, csrfMw(userMw.Apply(r)))

	return err
}