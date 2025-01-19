package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

// SQLite represents a connection to an SQLite database.
type SQLite struct {
	db *sql.DB
}

// NewSQLite creates a new connection to an SQLite database.
func NewSQLite(filePath string) (*SQLite, error) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}
	return &SQLite{db: db}, nil
}

// Close closes the connection to the SQLite database.
func (s *SQLite) Close() error {
	return s.db.Close()
}

// (Add more methods for database interactions as needed)
