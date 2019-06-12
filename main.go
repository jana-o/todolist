package main

import (
	"code/toDoList/config"
	"code/toDoList/handlers"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

type ToDo struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdat"`
}

type ToDoList struct {
	ToDos []ToDo
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db := config.InitDB()
	config.MigrateDB(db)

	e.File("/", "public/index.html")
	e.Static("/uploads", "public/uploads")
	e.GET("/todos", handlers.GetTodos(db))
	e.POST("/todos", handlers.Create(db))
	e.DELETE("/todos/:id", handlers.Delete(db))

	e.Logger.Fatal(e.Start(":9000"))
}
