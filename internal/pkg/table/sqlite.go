package table

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Fetch attempts to fetch the CID of the
func Fetch(path string, id string, format string) (error, string) {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("[sqlite] %s doesn't exist, initializing...\n", path)
		InitializeTable(path)
		return err, ""
	}
	return nil, ""
}

// InitializeTable creates an initialized table at path given
func InitializeTable(path string) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
CREATE TABLE ipfs (
  id TEXT NOT NULL,
  format TEXT NOT NULL,
  location TEXT,
  PRIMARY KEY (id, format)
);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
