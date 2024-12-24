package sqlstore

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func PostgresMigration(path string, dbUser, dbPass, dbAddr, dbPort, dbName string) {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbAddr, dbPort, dbName)
	sourceURL := fmt.Sprintf("file://%s", path)
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
