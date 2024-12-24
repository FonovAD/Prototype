package main

import (
	"log"

	"github.com/FonovAD/Prototype/internal/api"
	sqlstore "github.com/FonovAD/Prototype/internal/store/SQLstore"
)

func main() {
	pgConf := api.PostgresConfig{
		DBUser: "postgres",
		DBPass: "postgres",
		DBAddr: "db",
		DBPort: "5432",
		DBName: "linkshortener",
	}
	path := "./migrations"
	sqlstore.PostgresMigration(path,
		pgConf.DBUser,
		pgConf.DBPass,
		pgConf.DBAddr,
		pgConf.DBPort,
		pgConf.DBName,
	)
	log.Print("The database migration was successful.")
	if err := api.Start("info", "0.0.0.0:80", "127.0.0.1:80", pgConf); err != nil {
		panic(err)
	}
}
