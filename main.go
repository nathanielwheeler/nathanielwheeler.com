package main

import (
	"net/http"

	"nathanielwheeler.com/controllers"

	"github.com/gorilla/mux"
)

func main() {
	staticC := controllers.NewStatic()
	subsC := controllers.NewSubscribers()

	router := mux.NewRouter()
	router.HandleFunc("/", staticC.Home).Methods("GET")
	router.HandleFunc("/resume", staticC.Resume).Methods("GET")
	router.HandleFunc("/subscribe", subsC.New).Methods("GET")
	router.HandleFunc("/subscribe", subsC.Create).Methods("POST")
	http.ListenAndServe(":3000", router)
}
