package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//open db
	db, _ := sql.Open("sqlite3", "./data.db")
	//create table
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE
		);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	//insert data
	_, err = db.Exec(
		"INSERT INTO users(id, name, email) VALUES(?, ?, ?)",
		1, "Alice", "aliceinwonderland@gmail.com",
	)

	if err != nil {
		log.Fatal(err)
	}

	//read data
	var (
		id    int
		name  string
		email string
	)

	err = db.QueryRow(
		"SELECT id, name, email FROM users WHERE id = ?",
		1,
	).Scan(&id, &name, &email)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(id, name, email)
}
