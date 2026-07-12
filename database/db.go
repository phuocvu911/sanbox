package database

//used for open connection to sqlite3 database
import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Open() (*sql.DB, error) {
	return sql.Open("sqlite3", "./data.db")
}
