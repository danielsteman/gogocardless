package gocardless

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type UserAgreementRequestPayload struct {
	InstitionId        string   `json:"institution_id"`
	MaxHistoricalDays  int      `json:"max_historical_days"`
	AccessValidForDays int      `json:"access_valid_for_days"`
	AccessScope        []string `json:"access_scope"`
}

type UserAgreement struct {
	ID                 string    `json:"id"`
	Created            time.Time `json:"created"`
	MaxHistoricalDays  int       `json:"max_historical_days"`
	AccessValidForDays int       `json:"access_valid_for_days"`
	AccessScope        []string  `json:"access_scope"`
	Accepted           string    `json:"accepted"`
	InstitutionID      string    `json:"institution_id"`
}

func GetEndUserAgreement(institutionID string) (UserAgreement, error) {
	url := "https://bankaccountdata.gocardless.com/api/v2/agreements/enduser/"

	payload := UserAgreementRequestPayload{
		InstitionId:        institutionID,
		MaxHistoricalDays:  180,
		AccessValidForDays: 180,
		AccessScope:        []string{"balances", "details", "transactions"},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("Error marshalling JSON:", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	token, err := GetOrRefreshToken()
	if err != nil {
		log.Fatal("failed to get token:", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Access)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("failed to get user agreement:", err)
	}

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return UserAgreement{}, fmt.Errorf("failed to get user agreement: status code %d, response: %s", resp.StatusCode, string(body))
	}

	defer resp.Body.Close()
	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("failed to read response body:", err)
	}
	var userAgreement UserAgreement
	err = json.Unmarshal([]byte(jsonData), &userAgreement)
	if err != nil {
		log.Fatal("failed to get user agreement:", err)
	}

	return userAgreement, nil
}
