package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title string `json:"title"`
}

var db *gorm.DB

func main() {
	var err error
	// connect to database
	db, err = gorm.Open(postgres.Open("postgres://postgres:password@db/todos?sslmode=disable"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// run database migrations
	err = db.AutoMigrate(&Todo{})
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/todos", handleTodos).Methods("GET", "POST")

	log.Fatal(http.ListenAndServe(":8080", router))
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
