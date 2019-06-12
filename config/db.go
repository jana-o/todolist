package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

///Applications/Postgres.app/Contents/versions/latest/bin/psql -p5432

var DB *sql.DB

func InitDB() *sql.DB {
	var err error
	DB, err = sql.Open("postgres", "postgres://jana:password@localhost/todosdb?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Connected to todos database.")
	return DB
}

func MigrateDB(db *sql.DB) {
	sql := `
        CREATE TABLE IF NOT EXISTS todos(
			ID SERIAL PRIMARY KEY,
			TEXT           TEXT    NOT NULL,
			CREATEDAT		timestamp not null default now()
		);`

	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
}

//serial is postgres equivalent to autoincrement
