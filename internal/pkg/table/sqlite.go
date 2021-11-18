package table

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // Side effects
)

// OpenDB the sqlite database given by the filepath
func OpenDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, errors.New("db nil")
	}
	fmt.Printf("[sqlite] Opening '%s'\n", filepath)
	return db, nil
}

// Fetch attempts to fetch the CID of the file given by id and string
func Fetch(db *sql.DB, id string, format string) (string, error) {
	rows, err := db.Query(fmt.Sprintf(`
	SELECT location FROM ipfs
	WHERE id='%s' AND format='%s'
	`, id, format))
	if err != nil {
		return "", err
	}

	defer rows.Close()

	var out string
	for rows.Next() {
		err := rows.Scan(&out)
		if err != nil {
			return "", err
		}
	}
	fmt.Printf("[sqlite] fetched \"%s\" via \"%s\", \"%s\"\n", out, id, format)
	return out, nil
}

// Store attempts to insert an ipfs address pointing to a file with keys being id and format given
func Store(db *sql.DB, id string, format string, location string) error {
	fmt.Printf("[sqlite] attempting to insert (\"%s\", \"%s\", \"%s\")...\n", id, format, location)
	_, err := db.Exec(fmt.Sprintf(`
	INSERT INTO ipfs
	VALUES ("%s", "%s", "%s");`, id, format, location))
	return err
}

// InitializeTable creates an initialized table at path given
func InitializeTable(db *sql.DB) error {
	fmt.Printf("[sqlite] database doesn't exist, initializing...\n")
	_, err := db.Exec(`
	CREATE TABLE ipfs (
		id TEXT NOT NULL,
		format TEXT NOT NULL,
		location TEXT,
		PRIMARY KEY (id, format)
	);`)
	return err
}
