package handlers

import (
	"code/toDoList/todos"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type H map[string]interface{}

func GetTodos(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, todos.GetTodos(db))
	}
}

func Create(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		fmt.Println("enter Createhandler")
		var td todos.ToDo

		c.Bind(&td)

		text := c.FormValue("text")

		return c.JSON(http.StatusCreated, todos.Create(db, text))

	}
}

func Delete(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		err := todos.Delete(db, id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, H{"Deleted": id})
	}
}
