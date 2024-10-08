package main

import "log"

func main() {

	config := config{
		addr: ":8080",
	}

	app := &application{
		config,	
	}

	// instantiate the handler
	mux := app.mount()

	log.Fatal(app.run(mux))


}