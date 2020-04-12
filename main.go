package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func home(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-Type", "text/html")
	fmt.Fprint(res, "<h1>Nathaniel Wheeler</h1>")
}

func contact(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-Type", "text/html")
	fmt.Fprint(res,
		"To get in touch, please send an email to <a href=\"mailto:contact@nathanielwheeler.com\">contact@nathanielwheeler.com</a>.")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", home)
	router.HandleFunc("/contact", contact)
	http.ListenAndServe(":3000", router)
}
