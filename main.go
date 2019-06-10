package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

///Applications/Postgres.app/Contents/versions/latest/bin/psql -p5432
//grant all privileges on database todosdb to jana;
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

type ToDo struct {
	ID        int64
	Text      string
	CreatedAt time.Time
}

func main() {

	migrateDB(db)

	http.HandleFunc("/", index)
	http.HandleFunc("/todos", handleTodos)

	// http.HandleFunc("/todos/:id", getTodo)
	// http.HandleFunc("/todos/:id", updateTodo)
	http.HandleFunc("/todos/create", createForm)
	http.HandleFunc("/todos/delete", deleteTodo)
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

func createForm(w http.ResponseWriter, r *http.Request) {
	//get form values
	td := ToDo{}
	td.Text = r.FormValue("text")

	if td.Text == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	//insert values in db
	_, err := db.Exec(`INSERT INTO "todos" (text, createdat) VALUES ($1, $2)`, td.Text, time.Now())
	if err != nil {
		panic(err)
	}

	// fmt.Println(td)
	http.Redirect(w, r, "/todos", http.StatusSeeOther)

}

// func getTodo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	id := p.ByName("id")

// 	// Marshal into JSON
// 	uj, err := json.Marshal(td)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// Write content-type, statuscode, payload
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK) // 200
// 	fmt.Fprintf(w, "%s\n", uj)
// }

// func updateTodo(w http.ResponseWriter, r *http.Request) {
// 	// id := parse
// }

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	// id := p.ByName("id")
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "not found", 404)
		return
	}

	_, err := db.Exec("DELETE FROM todos WHERE id=$1", id)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/todos", http.StatusSeeOther)

}

func migrateDB(db *sql.DB) {
	sql := `
        CREATE TABLE IF NOT EXISTS todos(
			ID SERIAL PRIMARY KEY,
			TEXT           TEXT    NOT NULL,
			CREATEDAT		timestamptz not null default now()
		);`

	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
}

//timestamptz not null
//insert into todos (id, text, createdat) values (1, 't1', '1.1');
