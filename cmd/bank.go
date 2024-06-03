package main

import (
	"encoding/json"
	"net/http"

	"github.com/danielsteman/gogocardless/gocardless"
	"github.com/go-chi/chi/v5"
)

type bankResource struct{}

func (rs bankResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/list", listBanksHandler)

	return r
}

func listBanksHandler(w http.ResponseWriter, r *http.Request) {
	banks, err := gocardless.ListBanks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(banks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
