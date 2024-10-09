package main

import (
	"log"

	"github.com/salvatoreolivieri/go-api/internal/db"
	"github.com/salvatoreolivieri/go-api/internal/env"
	"github.com/salvatoreolivieri/go-api/internal/store"
)

func main() {

	config := config{
		addr: env.GetString("ADDR", ":8000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		config.db.addr,
		config.db.maxOpenConns,
		config.db.maxIdleConns,
		config.db.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("database connection pool established")

	store := store.NewStorage(db)

	app := &application{
		config,
		store,
	}

	// instantiate the handler
	mux := app.mount()

	log.Fatal(app.run(mux))

}
