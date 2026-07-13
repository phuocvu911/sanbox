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

	//set pragma
	if _, err := db.Exec(`PRAGMA journal_mode=WAL;`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(`PRAGMA foreign_keys=ON;`); err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	//create tables
	queryUsers := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE
		);
	`
	_, err = db.Exec(queryUsers)
	if err != nil {
		log.Fatal(err)
	}

	queryAccount := `
	CREATE TABLE IF NOT EXISTS accounts (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	owner TEXT NOT NULL,
    	balance INTEGER NOT NULL
	);
	`

	_, err = db.Exec(queryAccount)
	if err != nil {
		log.Fatal(err)
	}

	//insert data
	addUser(db, "alice", "aliceemial@gmail.com")
	addUser(db, "beta", "betal@gmail.com")
	addUser(db, "cena", "cenal@gmail.com")
	addUser(db, "dela", "deal@gmail.com")

	fmt.Println(listAccounts(db))
	err = transfer(db, 3, 4, 130)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(listAccounts(db))
}

type Account struct {
	Id      int
	Owner   string
	Balance int
}

func addAccount(db *sql.DB, name string, balance int) error {
	if name == "" {
		return fmt.Errorf("Name is empty")
	}
	_, err := db.Exec(
		"INSERT INTO accounts(owner, balance) VALUES(?,?)",
		name, balance,
	)
	if err != nil {
		return err
	}
	return nil
}

func listAccounts(db *sql.DB) ([]Account, error) {
	res := []Account{}
	rows, err := db.Query(`SELECT id, owner, balance FROM accounts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //
	for rows.Next() {
		var a Account
		err = rows.Scan(&a.Id, &a.Owner, &a.Balance)
		if err != nil {
			return res, err
		}
		res = append(res, a)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}
	return res, nil
}

func transfer(db *sql.DB, fromID, toID int, amount int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	var balanceFrom int
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = ?", fromID).Scan(&balanceFrom)
	if err != nil {
		return err
	}
	if balanceFrom < amount {
		tx.Rollback()
		return fmt.Errorf("insufficient funds: have %d, need %d", balanceFrom, amount)
	}
	_, err = tx.Exec(
		"UPDATE accounts SET balance = balance - ? WHERE id = ?",
		amount,
		fromID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		"UPDATE accounts SET balance = balance + ? WHERE id = ?",
		amount,
		toID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
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
