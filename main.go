package main

import (
	"net/http"

	"nathanielwheeler.com/views"

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
	contactView,
	subscribeView *views.View
)

// #region Handlers

func home(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-Type", "text/html")
	must(homeView.Render(res, nil))
}

func contact(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-Type", "text/html")
	must(contactView.Render(res, nil))
}

func subscribe(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	must(subscribeView.Render(res, nil))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// #endregion

func main() {
	homeView = views.NewView("app", "views/home.html")
	contactView = views.NewView("app", "views/contact.html")
	subscribeView = views.NewView("app", "views/subscribe.html")

	router := mux.NewRouter()
	router.HandleFunc("/", home)
	router.HandleFunc("/contact", contact)
	router.HandleFunc("/subscribe", subscribe)
	http.ListenAndServe(":3000", router)
}
