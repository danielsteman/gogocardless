package main

import (
	"encoding/json"
	"net/http"

	"github.com/danielsteman/gogocardless/gocardless"
	"github.com/go-chi/chi/v5"
)

type userResource struct{}

func (rs userResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/redirect", userRedirectHandler)

	return r
}

func userRedirectHandler(w http.ResponseWriter, r *http.Request) {
	institutionID := r.URL.Query().Get("institutionId")
	if institutionID == "" {
		http.Error(w, "institutionId query parameter is required", http.StatusBadRequest)
		return
	}

	redirectInfo, err := gocardless.GetEndUserRequisitionLink(institutionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(redirectInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
