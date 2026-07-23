package database

//create tables
import (
	"database/sql"
)

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE
		);
	`)
	if err != nil {
		return err
	}
	return nil
}
