package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Initialize sets up the SQLite database
func Initialize(path string) error {
	var err error
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	// Create tables
	query := `
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		path TEXT NOT NULL UNIQUE,
		size INTEGER,
		modified_time DATETIME,
		hash TEXT
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// Close closes the database connection
func Close() {
	if db != nil {
		db.Close()
	}
}
