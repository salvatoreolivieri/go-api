package main

import (
	"log"

	"github.com/salvatoreolivieri/go-api/internal/db"
	"github.com/salvatoreolivieri/go-api/internal/env"
	"github.com/salvatoreolivieri/go-api/internal/store"
)

func main() {
	addr := "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable"
	maxOpenConns := env.GetInt("DB_MAX_OPEN_CONNS", 3)
	maxIdleConns := env.GetInt("DB_MAX_IDLE_CONNS", 3)
	maxIdleTime := env.GetString("DB_MAX_IDLE_TIME", "15m")

	conn, err := db.New(addr, maxOpenConns, maxIdleConns, maxIdleTime)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
