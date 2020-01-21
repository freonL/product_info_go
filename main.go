package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/freonL/product_info_go/resources"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome home!")
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	port, _ := os.LookupEnv("PORT")
	fmt.Println("Starting the application at port :" + port)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homeLink)

	resources.AddRoutes(router.PathPrefix("/products").Subrouter())

	log.Fatal(http.ListenAndServe(":"+port, router))
}
