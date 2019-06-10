package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/lib/pq"
)

///Applications/Postgres.app/Contents/versions/latest/bin/psql -p5432
var db *sql.DB
var tpl *template.Template

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://jana:password@localhost/todosdb?sslmode=disable")
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")
	tpl = template.Must(template.ParseGlob("templates/*"))
}

//ToDO struct
type ToDo struct {
	ID        int64
	Text      string
	CreatedAt string
}

// CreatedAt time.Time `json:"createdAt"`

func main() {

	migrateDB(db)

	http.HandleFunc("/", index)
	http.HandleFunc("/todos", handleTodos)
	// http.HandleFunc("/todos/create", createForm)
	// http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	http.ListenAndServe(":9000", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func handleTodos(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		rows, err := db.Query("SELECT * FROM todos")
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		defer rows.Close()

		todos := make([]ToDo, 0)
		for rows.Next() {
			td := ToDo{}
			err := rows.Scan(&td.ID, &td.Text, &td.CreatedAt)
			if err != nil {
				fmt.Println("error GET db.Query rows.Scan")
			}
			todos = append(todos, td)
		}
		if err = rows.Err(); err != nil {
			panic(err)
		}
		fmt.Println(todos)
		tpl.ExecuteTemplate(w, "index.html", todos)

	}
}

// func createForm(w http.ResponseWriter, r *http.Request) {
// 	err := tpl.ExecuteTemplate(w, "create.html", nil)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 	}
// }

func migrateDB(db *sql.DB) {
	sql := `
        CREATE TABLE IF NOT EXISTS todos(
			ID INT PRIMARY KEY     NOT NULL,
			TEXT           TEXT    NOT NULL,
			CREATEDAT		TEXT NOT NULL
		);`

	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
}
