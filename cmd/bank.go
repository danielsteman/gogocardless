package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/danielsteman/gogocardless/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
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
	token, err := GetToken()

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
	}

	tokenResponse := TokenResponse{
		Token:   token,
		Elapsed: 0,
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, &tokenResponse)
}

type Credentials struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}

type Token struct {
	gorm.Model
	Access         string `json:"access"`
	AccessExpires  int    `json:"access_expires"`
	Refresh        string `json:"refresh"`
	RefreshExpires int    `json:"refresh_expires"`
}

type TokenResponse struct {
	*Token

	Elapsed int64 `json:"elapsed"`
}

func (rd *TokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}

func NewTokenResponse(token *Token) *TokenResponse {
	resp := &TokenResponse{Token: token}

	return resp
}

func GetToken() (*Token, error) {
	config, _ := config.LoadAppConfig(".env")
	url := "https://bankaccountdata.gocardless.com/api/v2/token/new/"
	credentials := Credentials{
		SecretID:  config.SecretID,
		SecretKey: config.SecretKey,
	}
	credentailsData, err := json.Marshal(credentials)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(credentailsData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var token Token
	if err := json.Unmarshal(responseBody, &token); err != nil {
		return nil, fmt.Errorf("error decoding token: %v", err)
	}

	return &token, nil
}

// func GetOrRefreshToken() (*Token, error) {

// }
