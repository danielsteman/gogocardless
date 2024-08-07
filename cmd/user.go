package main

import (
	"encoding/json"
	"net/http"

	"github.com/danielsteman/gogocardless/auth"
	"github.com/danielsteman/gogocardless/gocardless"
	"github.com/go-chi/chi/v5"
)

type userResource struct{}

func (rs userResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/redirect", userRedirectHandler)
	r.Get("/accounts", userAccountsHandler)
	r.Get("/accounts/{id}", userAccountTransactionsHandler)
	return r
}

func userRedirectHandler(w http.ResponseWriter, r *http.Request) {
	institutionID := r.URL.Query().Get("institutionId")
	if institutionID == "" {
		http.Error(w, "institutionId query parameter is required", http.StatusBadRequest)
		return
	}

	user := auth.GetUserFromContext(r)

	redirectInfo, err := gocardless.GetEndUserRequisitionLink(institutionID, user.Email)
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

	user := auth.GetUserFromContext(r)

	accountInfo, err := gocardless.GetEndUserAccountInfo(agreementRef, user.Email)
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

func userAccountTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "id")
	if accountId == "" {
		http.Error(w, "accountId path parameter is required", http.StatusBadRequest)
		return
	}

	user := auth.GetUserFromContext(r)

	transactions, err := gocardless.GetEndUserTransactions(accountId, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(transactions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
