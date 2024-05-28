package main

import (
	"encoding/json"
	"fmt"
	"io"
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

func (b *Bank) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (b *Bank) Bind(r *http.Request) error {
	return nil
}

func ListBanks(w http.ResponseWriter, r *http.Request) {
	token, err := gocardless.GetOrRefreshToken()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get token: %v", err), http.StatusInternalServerError)
		return
	}

	url := "https://bankaccountdata.gocardless.com/api/v2/institutions/?country=nl"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create request: %v", err), http.StatusInternalServerError)
		return
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Access)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get banks: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read response body: %v", err), http.StatusInternalServerError)
		return
	}

	var banks []Bank
	err = json.Unmarshal([]byte(jsonData), &banks)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse json: %v", err), http.StatusInternalServerError)
		return
	}

	bankList := make([]render.Renderer, len(banks))
	for i, bank := range banks {
		bankList[i] = &bank
	}

	render.Status(r, http.StatusOK)
	render.RenderList(w, r, bankList)
}
