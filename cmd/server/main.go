package main

import (
	"io"
	"log"
	"net/http"

	"github.com/Luc1808/todo-prj/internal/db"
	"github.com/Luc1808/todo-prj/internal/handlers"
	"github.com/Luc1808/todo-prj/internal/middlewares"
)

func main() {
	db.InitDB()

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("POST /register", handlers.RegisterHandler)
	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("GET /users", handlers.GetUsers)

	// Protected routes
	mux.Handle("GET /protected", middlewares.Authentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Welcome to the protected page!")
	})))

	port := ":8080"
	log.Printf("Server is running on port %s", port)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}
