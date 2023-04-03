package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var db *gorm.DB

func main() {
	var err error
	// connect to database
	db, err = gorm.Open(postgres.Open("postgres://postgres:password@db/todos?sslmode=disable"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database Connected")

	// run database migrations
	err = db.AutoMigrate(&Todo{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database Migrated")

	router := mux.NewRouter()

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Origin", "Accept"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
	router.Use(handlers.CORS(headers, origins, methods))
	log.Println("Applied CORS")

	router.Use(logRequest)
	router.HandleFunc("/todos", handleTodos).Methods("GET", "POST", "OPTIONS")
	log.Println("Applied Handlers")
	log.Println("Ready to Serve")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func logRequest(next http.Handler) http.Handler {

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func handleTodos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getTodos(w)
	case "POST":
		createTodo(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func getTodos(w http.ResponseWriter) {

	// retrieve all todos
	var todos []Todo
	db.Find(&todos)

	// send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {

	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// save new todo in db
	db.Create(&todo)

	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}
