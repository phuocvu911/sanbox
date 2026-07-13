package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//open db
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() //
	//create table
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE
		);
	`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	//insert data
	addUser(db, "alice", "aliceemial@gmail.com")
	addUser(db, "beta", "betal@gmail.com")
	addUser(db, "cena", "cenal@gmail.com")
	addUser(db, "dela", "deal@gmail.com")
	//read data
	fmt.Println(listUsers(db))
	fmt.Println(getUserByEmail(db, "deal@gmail.com"))
}

func addUser(db *sql.DB, name, email string) error {
	_, err := db.Exec(
		"INSERT INTO users(name, email) VALUES(?, ?)",
		name, email,
	)

	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) &&
			sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("UNIQUE contraint violation")
		} else {
			return err
		}
	}
	return nil
}

type User struct {
	ID    int
	Name  string
	Email string
}

func listUsers(db *sql.DB) ([]User, error) {
	res := []User{}
	rows, err := db.Query(`SELECT id, name, email FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //
	for rows.Next() {
		var u User
		err = rows.Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			return res, err
		}
		res = append(res, u)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}
	return res, nil
}

func getUserByEmail(db *sql.DB, email string) (*User, error) {
	var u User
	row := db.QueryRow(`SELECT id, name, email FROM users WHERE email = ?`, email)
	if err := row.Scan(&u.ID, &u.Name, &u.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
