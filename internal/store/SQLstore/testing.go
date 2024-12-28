package sqlstore

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

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

func SetupTestDB(t *testing.T, tokenForAdmin string) (*sql.DB, func(...string)) {
	t.Helper()
	databasePath := "./test"
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	schema := `PRAGMA foreign_keys = ON;
CREATE TABLE IF NOT EXISTS users(
UID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
Token TEXT NOT NULL,
Role varchar(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS links(
UID INTEGER REFERENCES users(UID) ON DELETE CASCADE,
OriginLink TEXT NOT NULL,
ShortLink TEXT UNIQUE NOT NULL,
CreatedAt integer,
ExpirationTime integer NOT NULL,
Status varchar(10) NOT NULL,
ScheduledDeletionTime integer NOT NULL,
PRIMARY KEY (UID, OriginLink)
);

INSERT INTO users(Token, Role) VALUES("%s", "admin");`

	if _, err := db.Exec(fmt.Sprintf(schema, tokenForAdmin)); err != nil {
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
		err := os.Remove(databasePath)
		if err != nil {
			fmt.Println(err)
		}
	}
}
