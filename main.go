package main

import (
	"context"
	"peltoar/v2/usercase"

	// "database/sql"
	// "encoding/json"
	// "fmt"
	"log"
	"net/http"
	"os"

	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/joho/godotenv"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	// "github.com/google/uuid"
)

var mongoClient *mongo.Client

func init() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("env load error", err)
	}

	log.Println("env file loaded")

	mongoClient, err = mongo.Connect(context.Background(),
		options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err != nil {
		log.Fatal("connection error", err)
	}

	err = mongoClient.Ping(context.Background(), readpref.Primary())

	if err != nil {
		log.Fatal("ping failed", err)
	}

	log.Println("mongo connected")
}

// type Article struct {
// 	Title   string `json:"Title"`
// 	Desc    string `json:"desc"`
// 	Content string `json:"content"`
// }

// type Articles []Article

// func allArticles(w http.ResponseWriter, r *http.Request) {
// 	articles := Articles{
// 		Article{Title: "Test Title", Desc: "Test Desc", Content: "Hello World"},
// 	}
// 	fmt.Printf("Endpoint Hit! All Articles Endpoint")
// 	json.NewEncoder(w).Encode(articles)
// }

// func testPostArticles(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Test POST Endpoint worked")
// }

// func homePage(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Homepage Endpoint Hit")
// }

// func handleRequests() {
// 	myRouter := mux.NewRouter().StrictSlash(true)
// 	myRouter.HandleFunc("/", homePage)
// 	myRouter.HandleFunc("/articles", allArticles).Methods("GET")
// 	myRouter.HandleFunc("/articles", testPostArticles).Methods("POST")
// 	log.Fatal(http.ListenAndServe(":8081", myRouter))
// }

// func sqlConfig() {
// 	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/testdb")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	defer db.Close()
// 	fmt.Println("Successfully connected to MySQL Database")

// 	insert, err := db.Query("INSERT INTO users VALUES('ELLIOT')")

// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	defer insert.Close()

// 	fmt.Println("Successfully inserted into users table")
// }

func main() {
	// handleRequests()
	defer mongoClient.Disconnect(context.Background())

	coll := mongoClient.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLLECTION_NAME"))
	// create employee service
	empService := usercase.EmployeeService{MongoCollection: coll}
	r := mux.NewRouter()

	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)
	r.HandleFunc("/employee", empService.CreateEmployee).Methods(http.MethodPost)
	r.HandleFunc("/employee/{id}", empService.GetEmployeeByID).Methods(http.MethodGet)
	r.HandleFunc("/employee", empService.GetAllEmployees).Methods(http.MethodGet)
	r.HandleFunc("/employee/{id}", empService.UpdateEmployeeByID).Methods(http.MethodPut)
	r.HandleFunc("/employee/{id}", empService.DeleteByEmployeeID).Methods(http.MethodDelete)
	r.HandleFunc("/employee", empService.DeleteAllEmployees).Methods(http.MethodDelete)

	log.Println("server is running on 4444")
	http.ListenAndServe(":4444", r)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("running..."))
}
