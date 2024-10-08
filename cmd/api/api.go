package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
}

type config struct {
    addr string
	// db
	// rateLimiter
}

func (app *application) mount() http.Handler{

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	
	r.Route("/v1", func(r chi.Router) {
		// endpoints
		r.Get("/health", app.healthCheckHandler)

	})

	return r
}

func (app  *application) run(mux http.Handler) error {


	server := &http.Server{
		Addr: app.config.addr, 
		Handler: mux, // route handler
		WriteTimeout: time.Second * 30, // maximum duration before timouts write for the response
		ReadTimeout: time.Second * 10, // maximum duration before timing out read for the request
		IdleTimeout: time.Minute, // maximum duration before timing out idle connections
	}

	log.Printf("listening on %s", app.config.addr)

	return server.ListenAndServe()
}