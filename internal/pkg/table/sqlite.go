package table

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type FileItem struct {
	Path   string
	Id     string
	Format string
}

func OpenDB(filepath string) (error, *sql.DB) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return err, nil
	}
	if db == nil {
		return errors.New("db nil"), nil
	}
	fmt.Printf("[sqlite] Opening '%s'\n", filepath)
	return nil, db
}

// Fetch attempts to fetch the CID of the file given by id and string
func Fetch(db *sql.DB, id string, format string) (error, string) {
	rows, err := db.Query(fmt.Sprintf(`
	SELECT location FROM ipfs
	WHERE id='%s' AND format='%s'
	`, id, format))
	if err != nil {
		return err, ""
	}

	defer rows.Close()

	var out string
	for rows.Next() {
		err := rows.Scan(&out)
		if err != nil {
			return err, ""
		}
	}
	fmt.Printf("[sqlite] fetched \"%s\" via \"%s\", \"%s\"\n", out, id, format)
	return nil, out
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
