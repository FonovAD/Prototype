package main

import (
	"flag"
	"log"

	"github.com/FonovAD/Prototype/config"
	"github.com/FonovAD/Prototype/internal/api"
	sqlstore "github.com/FonovAD/Prototype/internal/store/SQLstore"
)

func main() {
	UseSQLite3 := flag.Bool("sql", true, "true - use SQLite3, false - user PostgreSQL")
	flag.Parse()

	configPath := "./config/default.yaml"

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}
	path := "./migrations"
	if !*UseSQLite3 {
		sqlstore.PostgresMigration(path,
			cfg.Postgres.User,
			cfg.Postgres.Password,
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.Database,
		)
		log.Print("The database migration was successful.")
	}
	if err := api.Start(cfg, *UseSQLite3); err != nil {
		panic(err)
	}
}
