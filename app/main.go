// @title           ToDo API
// @version         1.0
// @description     A simple ToDo API built with Go and PostgreSQL.
// @termsOfService  http://example.com/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey  BearerToken
// @in                          header
// @name                        Authorization

package main

import (
	"io"
	"log"
	"net/http"

	_ "github.com/Luc1808/todo-prj/docs"
	"github.com/Luc1808/todo-prj/internal/db"
	"github.com/Luc1808/todo-prj/internal/handlers"
	"github.com/Luc1808/todo-prj/internal/middlewares"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	db.InitDB()

	mux := http.NewServeMux()

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Public routes
	mux.HandleFunc("POST /register", handlers.RegisterHandler)
	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("POST /refresh", handlers.RefreshHandler)
	// mux.HandleFunc("GET /users", handlers.GetUsers)

	// Protected routes
	mux.Handle("POST /tasks", middlewares.Authentication(http.HandlerFunc(handlers.CreateTodoHandler)))
	mux.Handle("GET /tasks/{id}", middlewares.Authentication(http.HandlerFunc(handlers.GetTodoByID)))
	mux.Handle("PUT /tasks/{id}", middlewares.Authentication(http.HandlerFunc(handlers.UpdateTodo)))
	mux.Handle("DELETE /tasks/{id}", middlewares.Authentication(http.HandlerFunc(handlers.DeleteTodo)))
	mux.Handle("GET /tasks", middlewares.Authentication(http.HandlerFunc(handlers.GetAllTodosWithPagination)))
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
