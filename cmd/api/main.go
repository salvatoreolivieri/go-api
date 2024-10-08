package main

import (
	"log"

	"github.com/salvatoreolivieri/go-api/internal/env"
)

func main() {

	config := config{
		addr: env.GetString("ADDR", "8080"),
	}

	app := &application{
		config,	
	}

	// instantiate the handler
	mux := app.mount()

	log.Fatal(app.run(mux))


}