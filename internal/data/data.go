package data

import (
	"database/sql"
	"log"
	"os"

	// SQLite3 Driver
	_ "github.com/mattn/go-sqlite3"
)

// DB contains a pointer to a pool and methods for interfacing with
// the db.
type DB struct {
	pool *sql.DB
}

// Connect opens the pool connection
func (db *DB) Connect() {
	pool, err := sql.Open("sqlite3", os.Getenv("SQLITE_FILE"))
	if err != nil {
		log.Println(err)
	}
	db.pool = pool
}

// Disconnect closes the pool connection
func (db *DB) Disconnect() {
	db.pool.Close()
}

// Setup executes a query to create tables
func (db *DB) Setup() (err error) {
	q := `	DROP TABLE IF EXISTS links;
		CREATE TABLE links (
		id SERIAL PRIMARY KEY,
		target_url TEXT,
		short_url TEXT
		);`
	_, err = db.pool.Exec(q)
	if err != nil {
		log.Fatal(err)
	}
	return
}
