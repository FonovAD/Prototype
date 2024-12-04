package main

import "github.com/FonovAD/Prototype/internal/api"

func main() {
	if err := api.Start("info", "127.0.0.1:80"); err != nil {
		panic(err)
	}
}
