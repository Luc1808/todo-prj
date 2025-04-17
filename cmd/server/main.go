package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Luc1808/todo-prj/internal/db"
	"github.com/Luc1808/todo-prj/internal/handlers"
)

func main() {

	db.InitDB()

	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {

		type RequestBody struct {
			Username string `json:"username"`
		}

		var user RequestBody
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		io.WriteString(w, user.Username)
	})
	http.HandleFunc("POST /register", handlers.RegisterHandler)
	http.HandleFunc("POST /login", handlers.LoginHandler)
	http.HandleFunc("GET /users", handlers.GetUsers)

	port := ":8080"
	log.Printf("Server is running on port %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
