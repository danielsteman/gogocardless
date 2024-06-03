package gocardless

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Bank struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	BIC                  string   `json:"bic"`
	TransactionTotalDays string   `json:"transaction_total_days"`
	Countries            []string `json:"countries"`
	Logo                 string   `json:"logo"`
}

func ListBanks() ([]Bank, error) {
	token, err := GetOrRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := "https://bankaccountdata.gocardless.com/api/v2/institutions/?country=nl"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Access)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get banks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get banks: status code %d, response: %s", resp.StatusCode, string(body))
	}

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var banks []Bank
	err = json.Unmarshal(jsonData, &banks)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}
	return banks, nil
}
