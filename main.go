package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Model (field name start with capital)
type product struct {
	ID       string  `json:"_id"`
	Title    string  `json:"title"`
	Category string  `json:"category"`
	Price    float32 `json:"price"`
	Pic_url  string  `json:"pic_url"`
}

type allProducts []product

// dummy data
var products = allProducts{
	{

		Title:    "Tomato",
		Category: "Fruit",
		Price:    3.5,
		Pic_url:  "",
		ID:       "1",
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome home!")
}

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	// fmt.Println( len( products) )

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct product
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newProduct)
	products = append(products, newProduct)
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newProduct)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	docID := mux.Vars(r)["id"]
	w.Header().Set("Content-Type", "application/json")
	for _, doc := range products {
		if doc.ID == docID {
			json.NewEncoder(w).Encode(doc)
		}
	}
}

func updateProduct(w http.ResponseWriter, r *http.Request) {

}

func deleteProduct(w http.ResponseWriter, r *http.Request) {

}

func main() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homeLink)
	router.HandleFunc("/products", getAllProducts).Methods("GET")
	router.HandleFunc("/products", createProduct).Methods("POST")
	router.HandleFunc("/products/{id}", getProduct).Methods("GET")
	router.HandleFunc("/products/{id}", updateProduct).Methods("PATCH")
	router.HandleFunc("/products/{id}", deleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
