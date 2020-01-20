package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/freonL/product_info_go/resources"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome home!")
}

func main() {

	fmt.Println("Starting the application...")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homeLink)

	resources.AddRoutes(router)

	port, _ := os.LookupEnv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
