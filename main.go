package main

import (
	"net/http"

	"nathanielwheeler.com/views"

	"github.com/gorilla/mux"
)

var homeView, contactView *views.View

func home(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-Type", "text/html")
	must(homeView.Render(res, nil))
}

func contact(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-Type", "text/html")
	must(contactView.Render(res, nil))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	homeView = views.NewView("app", "views/home.html")
	contactView = views.NewView("app", "views/contact.html")

	router := mux.NewRouter()
	router.HandleFunc("/", home)
	router.HandleFunc("/contact", contact)
	http.ListenAndServe(":3000", router)
}
