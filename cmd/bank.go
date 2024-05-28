package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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
	token, err := createNewToken()

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get token: %v", err), http.StatusInternalServerError)
		return
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
	rd.Elapsed = 10
	return nil
}

func createNewToken() (*Token, error) {
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

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var token Token
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("error decoding token: %v", err)
	}

	return &token, nil
}

// GetOrRefreshToken retrieves an existing token or generates a new one if necessary
func GetOrRefreshToken() (*Token, error) {
	db, err := getDB()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	var token Token
	result := db.Order("created_at desc").First(&token)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("error retrieving token: %w", result.Error)
	}

	if result.Error == nil {
		// Check if the token is still valid
		expiresAt := token.CreatedAt.Add(time.Duration(token.AccessExpires) * time.Second)
		if time.Now().Before(expiresAt) {
			return &token, nil
		}
	}

	// Token is expired or not found, create a new token
	newToken, err := createNewToken()
	if err != nil {
		return nil, fmt.Errorf("error creating new token: %w", err)
	}

	// Save the new token to the database
	if _, err := dbCreateToken(newToken); err != nil {
		return nil, fmt.Errorf("error saving new token: %w", err)
	}

	return newToken, nil
}

func dbCreateToken(token *Token) (string, error) {
	db, err := getDB()
	if err != nil {
		return "", fmt.Errorf("error connecting to the database: %w", err)
	}
	if err := db.Create(token).Error; err != nil {
		return "", fmt.Errorf("error creating token: %w", err)
	}
	return "Token created successfully", nil
}
