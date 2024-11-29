package sqlstore

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func TestDB(t *testing.T, databasePath string) (*sql.DB, func(...string)) {
	t.Helper()
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
	return db, func(tables ...string) {
		if len(tables) > 0 {
			_, err := db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE;", strings.Join(tables, ", ")))
			if err != nil {
				panic(fmt.Sprintf("the test table could not be cleared: %s", err))
			}
		}
		db.Close()
	}
}

func SetupTestDB(t *testing.T) (*sql.DB, func(...string)) {
	t.Helper()
	db, err := sql.Open("sqlite3", "./test")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	schema := `
CREATE TABLE users(
UID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
Token TEXT NOT NULL,
Role varchar(10) NOT NULL
);

CREATE TABLE links(
UID INTEGER REFERENCES users(UID) ON DELETE CASCADE,
OriginLink TEXT NOT NULL,
ShortLink TEXT,
CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ExpirationTime TIMESTAMP NOT NULL,
Status varchar(10) NOT NULL,
ScheduledDeletionTime TIMESTAMP
);`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Failed to setup test database schema: %v", err)
	}
	return db, func(tables ...string) {
		if len(tables) > 0 {
			for _, table := range tables {
				_, err := db.Exec("DELETE FROM " + table)
				if err != nil {
					log.Printf("Failed to clear table %s: %v", table, err)
				}
			}
		}
		db.Close()
	}
}
