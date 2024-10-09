package main

import (
	"log"

	"github.com/salvatoreolivieri/go-api/internal/env"
	"github.com/salvatoreolivieri/go-api/internal/store"
)

func main() {

	config := config{
		addr: env.GetString("ADDR", ":8000"),
	}

	store := store.NewStorage(nil)

	app := &application{
		config,
		store,
	}

	// instantiate the handler
	mux := app.mount()

	log.Fatal(app.run(mux))

}
