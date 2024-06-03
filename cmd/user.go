package main

import "github.com/go-chi/chi/v5"

type userResource struct{}

func (rs userResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Get("/redirect", ListBanks)

	return r
}
