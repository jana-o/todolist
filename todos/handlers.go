package todos

import (
	"code/toDoList/config"
	"database/sql"
	"fmt"
	"net/http"
)

func Index(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	tds, err := GetTodos()
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "index.html", tds)
}

func CreateHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	td, err := CreateTodo(req)
	if err != nil {
		panic(err)
	}
	config.TPL.ExecuteTemplate(w, "created.html", td)
	// http.Redirect(w, req, "/todos", http.StatusSeeOther)

}

func DeleteHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	err := DeleteTodo(req)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, req, "/todos", http.StatusSeeOther)
}

func UpdateHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	td, err := OneTodo(req)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, req)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "update.html", td)
}

func UpdatePostHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	td, err := UpdateTodo(req)
	if err != nil {
		fmt.Println("err UpdatePostHandler")
		http.Error(w, http.StatusText(406), http.StatusBadRequest)
		return
	}

	config.TPL.ExecuteTemplate(w, "updated.html", td)
}
