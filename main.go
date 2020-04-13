package main

import (
	"net/http"

	"nathanielwheeler.com/views"
	"nathanielwheeler.com/controllers"

	"github.com/gorilla/mux"
)

// #region TODO Page adding procedure:
/*

/views
- create view template (.html)
- define as "yield"

/views/layouts/navbar.html
- Add to navbar

main.go
- add view variable (*views.View)
- add handler func
- in main()...
	- initialize view
	- call handler from router

*/
// #endregion

var (
	homeView,
	resumeView *views.View
)

// #region Handlers

func home(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-Type", "text/html")
	must(homeView.Render(res, nil))
}

func resume(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-Type", "text/html")
	must(resumeView.Render(res, nil))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// #endregion

func main() {
	homeView = views.NewView("app", "views/home.html")
	resumeView = views.NewView("app", "views/resume.html")
	subsController := controllers.NewSubscribers()

	router := mux.NewRouter()
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/resume", resume).Methods("GET")
	router.HandleFunc("/subscribe", subsController.New).Methods("GET")
	router.HandleFunc("/subscribe", subsController.Create).Methods("POST")
	http.ListenAndServe(":3000", router)
}
