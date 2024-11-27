package main

import "github.com/FonovAD/Prototype/internal/api"

func main() {
	if err := api.Start("info", "127.0.0.1:8080"); err != nil {
		panic(err)
	}
}
