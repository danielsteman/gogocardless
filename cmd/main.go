package main

import (
	"log"
	"net/http"

	"github.com/danielsteman/gogocardless/db"
	"github.com/danielsteman/gogocardless/gocardless"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	db, err := db.GetDB("gogocardless")
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	err = db.AutoMigrate(&gocardless.Token{})
	if err != nil {
		panic("failed to migrate database")
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sup dawg!"))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Mount("/banks", bankResource{}.Routes())

	http.ListenAndServe(":3333", r)
}
