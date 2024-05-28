package main

import (
	"fmt"
	"net/http"

	"github.com/danielsteman/gogocardless/gocardless"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type bankResource struct{}

func (rs bankResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/list", ListBanks)

	return r
}

type Bank struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	BIC                  string   `json:"bic"`
	TransactionTotalDays string   `json:"transaction_total_days"`
	Countries            []string `json:"countries"`
	Logo                 string   `json:"logo"`
}

type ListBanksResponse struct {
	Banks   []*Bank
	Elapsed int64 `json:"elapsed"`
}

func ListBanks(w http.ResponseWriter, r *http.Request) {
	token, err := gocardless.GetOrRefreshToken()

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get token: %v", err), http.StatusInternalServerError)
		return
	}

	tokenResponse := gocardless.TokenResponse{
		Token:   token,
		Elapsed: 0,
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, &tokenResponse)
}
