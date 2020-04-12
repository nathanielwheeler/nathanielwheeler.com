package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "<h1>Le Gestionnaire de Photos</h1>")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)
}
