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
	r.Post("/redirect", userRedirectHandler)
	r.Get("/accounts", userAccountsHandler)

	return r
}

type RedirectRequest struct {
	InstitutionID string `json:"institutionId"`
	UserEmail     string `json:"userEmail"`
}

func userRedirectHandler(w http.ResponseWriter, r *http.Request) {
	var redirectRequest RedirectRequest
	err := json.NewDecoder(r.Body).Decode(&redirectRequest)
	if err != nil {
		http.Error(w, "error parsing redirect request body", http.StatusInternalServerError)
		return
	}

	if redirectRequest.InstitutionID == "" {
		http.Error(w, "InstitutionId is required in request body", http.StatusBadRequest)
		return
	}

	if redirectRequest.UserEmail == "" {
		http.Error(w, "UserEmail is required in request body", http.StatusBadRequest)
		return
	}

	redirectInfo, err := gocardless.GetEndUserRequisitionLink(redirectRequest.InstitutionID)
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

func userAccountsHandler(w http.ResponseWriter, r *http.Request) {
	agreementRef := r.URL.Query().Get("agreementRef")
	if agreementRef == "" {
		http.Error(w, "agreementRef query parameter is required", http.StatusBadRequest)
		return
	}

	accountInfo, err := gocardless.GetEndUserAccountInfo(agreementRef)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(accountInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
