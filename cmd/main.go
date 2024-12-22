package main

import (
	"github.com/FonovAD/Prototype/internal/api"
	sqlstore "github.com/FonovAD/Prototype/internal/store/SQLstore"
)

func main() {
	path := "./migrations"
	DB_user := "aleksandr"
	DB_pass := ""
	DB_addr := "localhost"
	DB_port := "5432"
	DB_name := "linkshortener"
	sqlstore.PostgresMigration(path, DB_user, DB_pass, DB_addr, DB_port, DB_name)

	if err := api.Start("info", "127.0.0.1:80", "127.0.0.1:80"); err != nil {
		panic(err)
	}
}
