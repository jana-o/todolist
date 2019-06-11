package main

import (
	"code/toDoList/config"
	"code/toDoList/todos"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {

	config.MigrateDB(config.DB)

	http.HandleFunc("/", index)
	http.HandleFunc("/todos", todos.Index)
	http.HandleFunc("/todos/create", todos.CreateHandler)
	http.HandleFunc("/todos/delete", todos.DeleteHandler)
	http.HandleFunc("/todos/update", todos.UpdateHandler)
	http.HandleFunc("/todos/update/post", todos.UpdatePostHandler)
	// http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	http.ListenAndServe(":9000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todos", http.StatusSeeOther)
}
