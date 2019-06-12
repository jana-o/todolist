package todos

import (
	"code/toDoList/config"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type ToDo struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdat"`
}

type ToDoList struct {
	ToDos []ToDo
}

func GetTodos(db *sql.DB) ToDoList {

	rows, err := db.Query("SELECT * FROM todos")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	todos := ToDoList{}

	for rows.Next() {
		td := ToDo{}
		err := rows.Scan(&td.ID, &td.Text, &td.CreatedAt)
		if err != nil {
			fmt.Println("error GET db.Query rows.Scan")
		}
		todos.ToDos = append(todos.ToDos, td)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
	fmt.Println(todos)
	return todos
}

func Create(db *sql.DB, text string) ToDo {

	fmt.Println("enter Create")

	row := db.QueryRow("INSERT INTO todos (text) VALUES ($1)", text)

	td := ToDo{}
	err := row.Scan(&td.ID, &td.Text, &td.CreatedAt)
	if err != nil {
		fmt.Println("error GET db.Query rows.Scan")
	}
	return td

}

func Delete(db *sql.DB, id string) error {

	_, err := config.DB.Exec("DELETE FROM todos WHERE id=$1;", id)
	if err != nil {
		return errors.New("500. Internal Server Error")
	}
	return nil
}
