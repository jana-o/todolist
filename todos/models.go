package todos

import (
	"code/toDoList/config"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type ToDo struct {
	ID        int
	Text      string
	CreatedAt time.Time
}

func GetTodos() ([]ToDo, error) {

	rows, err := config.DB.Query("SELECT * FROM todos")
	if err != nil {
		return nil, err
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
	return todos, nil
}

func OneTodo(req *http.Request) (ToDo, error) {

	id := req.FormValue("id")
	td := ToDo{}

	row := config.DB.QueryRow("SELECT*FROM todos WHERE id=$1", id)
	err := row.Scan(&td.ID, &td.Text, &td.CreatedAt)
	if err != nil {
		return td, nil
	}

	return td, nil
}

func UpdateTodo(req *http.Request) (ToDo, error) {

	id := req.FormValue("id")

	td := ToDo{}
	newid := req.FormValue("id")
	td.Text = req.FormValue("text")
	td.CreatedAt = time.Now()

	new, err := strconv.Atoi(newid)
	if err != nil {
		fmt.Println("error converting")
	}
	td.ID = new
	fmt.Println("enter UpdateTodo", td, new)

	// insert values
	_, err = config.DB.Exec("UPDATE todos SET id=$1, text=$2, createdat=$3 WHERE id=$4;", td.ID, td.Text, td.CreatedAt, id)
	if err != nil {
		fmt.Println("error updating db")
		return td, err
	}
	return td, nil
}

func CreateTodo(req *http.Request) (ToDo, error) {
	//get form values
	td := ToDo{}
	td.Text = req.FormValue("text")

	if td.Text == "" {
		return td, errors.New("400. Bad request. CreateToDo")
	}

	//insert values in db
	_, err := config.DB.Exec(`INSERT INTO "todos" (text, createdat) VALUES ($1, $2)`, td.Text, time.Now())
	if err != nil {
		panic(err)
	}

	return td, nil
}

func DeleteTodo(req *http.Request) error {
	id := req.FormValue("id")
	if id == "" {
		return errors.New("not found")
	}

	_, err := config.DB.Exec("DELETE FROM todos WHERE id=$1", id)
	if err != nil {
		return errors.New("500. Internal Server Error. DeleteTodo from db")
	}
	return nil
}
