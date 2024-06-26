package main

import (
	"log"
	"net/http"

	"github.com/danielsteman/gogocardless/auth"
	"github.com/danielsteman/gogocardless/config"
	"github.com/danielsteman/gogocardless/db"
	"github.com/danielsteman/gogocardless/gocardless"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	config.LoadAppConfig(".env")
	db, err := db.GetDB()
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	err = db.AutoMigrate(
		&gocardless.Token{},
		&gocardless.DBRequisition{},
		&gocardless.AccountInfo{},
		&gocardless.Account{},
		&gocardless.TransactionAmount{},
		&gocardless.BookedTransaction{},
		&gocardless.PendingTransaction{},
		&gocardless.Transactions{},
	)
	if err != nil {
		panic("failed to migrate database")
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(auth.VerifyToken)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sup dawg!"))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Mount("/api/banks", bankResource{}.Routes())
	r.Mount("/api/user", userResource{}.Routes())

	http.ListenAndServe(":3333", r)
}
