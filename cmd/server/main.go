package main

import (
	"io"
	"net/http"

	"github.com/Luc1808/todo-prj/internal/db"
)

func main() {

	db.InitDB()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "This is my website!")
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
