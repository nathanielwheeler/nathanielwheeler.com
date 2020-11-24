package server

func (s *server) routes() {
	r := s.router

		// Public Routes
		publicHandler := http.FileServer(http.Dir("./public/"))
		r.PathPrefix("/images/").
			Handler(publicHandler)
		r.PathPrefix("/stylesheets/").
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
}