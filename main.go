package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func homeLink( w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome home!")
}


func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	log.Fatal(http.ListenAndServe(":8080", router))
}