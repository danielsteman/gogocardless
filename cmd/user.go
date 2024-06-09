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
	r.Post("/accounts", userAccountsHandler)

	return r
}

type RedirectRequest struct {
	InstitutionID string `json:"institutionId"`
	UserEmail     string `json:"userEmail"`
}

type AccountsRequest struct {
	AgreementRef string `json:"agreementRef"`
	UserEmail    string `json:"userEmail"`
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

	redirectInfo, err := gocardless.GetEndUserRequisitionLink(redirectRequest.InstitutionID, redirectRequest.UserEmail)
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
	var accountsRequest AccountsRequest
	err := json.NewDecoder(r.Body).Decode(&accountsRequest)
	if err != nil {
		http.Error(w, "error parsing accounts request body", http.StatusBadRequest)
		return
	}

	accountInfo, err := gocardless.GetEndUserAccountInfo(accountsRequest.AgreementRef, accountsRequest.UserEmail)
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
