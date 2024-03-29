package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("postgres", "postgres://jana:password@localhost/todosdb?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")
}

func MigrateDB(db *sql.DB) {
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

//serial is postgres equivalent to autoincrement
