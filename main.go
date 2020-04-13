package main

import (
	"net/http"

	"nathanielwheeler.com/views"

	"github.com/gorilla/mux"
)

var homeView, contactView *views.View

func home(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-Type", "text/html")
	err := homeView.Template.ExecuteTemplate(res, homeView.Layout, nil)
	if err != nil {
		panic(err)
	}
}

func contact(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-Type", "text/html")
	err := contactView.Template.ExecuteTemplate(res, contactView.Layout, nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	homeView = views.NewView("app", "views/home.html")
	contactView = views.NewView("app", "views/contact.gohtml")

	router := mux.NewRouter()
	router.HandleFunc("/", home)
	router.HandleFunc("/contact", contact)
	http.ListenAndServe(":3000", router)
}
