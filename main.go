package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// Model (field name start with capital)
type Product struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title    *string            `json:"title,omitempty" bson:"title,omitempty"`
	Category *string            `json:"category,omitempty" bson:"category,omitempty"`
	Price    *float64           `json:"price,omitempty" bson:"price,omitempty"`
	PicURL   *string            `json:"pic_url,omitempty" bson:"pic_url,omitempty"`
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome home!")
}

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection := client.Database("product_info").Collection("products")

	// Find Multiple Documents
	findOptions := options.Find()
	// findOptions.SetLimit(10)

	// Here's an array in which you can store the decoded documents
	var results []*Product

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var doc Product
		err := cur.Decode(&doc)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &doc)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	json.NewEncoder(w).Encode(results)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection := client.Database("product_info").Collection("products")
	var newDoc Product

	_ = json.NewDecoder(r.Body).Decode(&newDoc)
	result, err := collection.InsertOne(context.TODO(), newDoc)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(result)

}

func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	paramId := mux.Vars(r)["id"]
	docID, _ := primitive.ObjectIDFromHex(paramId)

	var doc Product
	collection := client.Database("product_info").Collection("products")
	ctx, _ := context.WithTimeout(context.TODO(), 30*time.Second)

	filter := bson.D{{"_id", docID}}
	err := collection.FindOne(ctx, filter).Decode(&doc)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(doc)

}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	paramId := mux.Vars(r)["id"]
	docID, _ := primitive.ObjectIDFromHex(paramId)
	fmt.Println(docID)
	json.NewEncoder(w).Encode(json.NewDecoder(r.Body))
	var newDoc Product
	err := json.NewDecoder(r.Body).Decode(&newDoc)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("product_info").Collection("products")

	filter := bson.D{{"_id", docID}}
	update := bson.D{{"$set", newDoc}}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(updateResult)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	paramId := mux.Vars(r)["id"]
	docID, _ := primitive.ObjectIDFromHex(paramId)

	collection := client.Database("product_info").Collection("products")

	filter := bson.D{{"_id", docID}}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(deleteResult)

}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	mongoURI, _ := os.LookupEnv("MONGO_URI")

	fmt.Println("Starting the application...")
	// Setup Database connection
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, _ = mongo.Connect(ctx, clientOptions)

	// Routing endpoints
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homeLink)
	router.HandleFunc("/products", getAllProducts).Methods("GET")
	router.HandleFunc("/products", createProduct).Methods("POST")
	router.HandleFunc("/products/{id}", getProduct).Methods("GET")
	router.HandleFunc("/products/{id}", updateProduct).Methods("PATCH")
	router.HandleFunc("/products/{id}", deleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
