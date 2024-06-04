package gocardless

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/danielsteman/gogocardless/config"
	"github.com/google/uuid"
)

type UserAgreementRequestPayload struct {
	InstitutionId      string   `json:"institution_id"`
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

type RequisitionPayload struct {
	Redirect      string `json:"redirect"`
	InstitutionID string `json:"institution_id"`
	Reference     string `json:"reference"`
	Agreement     string `json:"agreement"`
	UserLanguage  string `json:"user_language"`
}

type RedirectInfo struct {
	ID           string   `json:"id"`
	Redirect     string   `json:"redirect"`
	Status       string   `json:"status"`
	Agreement    string   `json:"agreement"`
	Accounts     []string `json:"accounts"`
	Reference    string   `json:"reference"`
	UserLanguage string   `json:"user_language"`
	Link         string   `json:"link"`
}

type AccountInfo struct {
	ID         string   `json:"id"`
	Status     string   `json:"status"`
	Agreements string   `json:"agreements"`
	Accounts   []string `json:"accounts"`
	Reference  string   `json:"reference"`
}

func GetEndUserAgreement(institutionID string) (UserAgreement, error) {
	url := "https://bankaccountdata.gocardless.com/api/v2/agreements/enduser/"

	payload := UserAgreementRequestPayload{
		InstitutionId:      institutionID,
		MaxHistoricalDays:  180,
		AccessValidForDays: 180,
		AccessScope:        []string{"balances", "details", "transactions"},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return UserAgreement{}, fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return UserAgreement{}, fmt.Errorf("error creating request: %w", err)
	}

	token, err := GetOrRefreshToken()
	if err != nil {
		return UserAgreement{}, fmt.Errorf("failed to get token: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Access)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return UserAgreement{}, fmt.Errorf("failed to get user agreement: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return UserAgreement{}, fmt.Errorf("failed to get user agreement: status code %d, response: %s", resp.StatusCode, string(body))
	}

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserAgreement{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var userAgreement UserAgreement
	err = json.Unmarshal(jsonData, &userAgreement)
	if err != nil {
		return UserAgreement{}, fmt.Errorf("failed to unmarshal user agreement: %w", err)
	}

	return userAgreement, nil
}

func GetEndUserRequisitionLink(institutionID string) (RedirectInfo, error) {
	userAgreement, err := GetEndUserAgreement(institutionID)
	if err != nil {
		return RedirectInfo{}, fmt.Errorf("failed to get user agreement: %w", err)
	}

	newReference := uuid.New().String()
	url := "https://bankaccountdata.gocardless.com/api/v2/requisitions/"

	payload := RequisitionPayload{
		Redirect:      config.Config.RedirectURL,
		InstitutionID: userAgreement.InstitutionID,
		Reference:     newReference,
		Agreement:     userAgreement.ID,
		UserLanguage:  "EN",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return RedirectInfo{}, fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return RedirectInfo{}, fmt.Errorf("error creating request: %w", err)
	}

	token, err := GetOrRefreshToken()
	if err != nil {
		return RedirectInfo{}, fmt.Errorf("failed to get token: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Access)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return RedirectInfo{}, fmt.Errorf("failed to get redirect info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return RedirectInfo{}, fmt.Errorf("failed to get redirect info: status code %d, response: %s", resp.StatusCode, string(body))
	}

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return RedirectInfo{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var redirectInfo RedirectInfo
	err = json.Unmarshal(jsonData, &redirectInfo)
	if err != nil {
		return RedirectInfo{}, fmt.Errorf("failed to unmarshal redirect info: %w", err)
	}

	return redirectInfo, nil
}

func GetEndUserAccountInfo(agreementID string) (AccountInfo, error) {
	url := "https://bankaccountdata.gocardless.com/api/v2/requisitions/" + agreementID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return AccountInfo{}, fmt.Errorf("error creating request: %w", err)
	}

	token, err := GetOrRefreshToken()
	if err != nil {
		return AccountInfo{}, fmt.Errorf("failed to get token: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Access)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AccountInfo{}, fmt.Errorf("failed to get account info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return AccountInfo{}, fmt.Errorf("failed to get account info: status code %d, response: %s", resp.StatusCode, string(body))
	}

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return AccountInfo{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var accountInfo AccountInfo
	err = json.Unmarshal(jsonData, &accountInfo)
	if err != nil {
		return AccountInfo{}, fmt.Errorf("failed to unmarshal account info: %w", err)
	}

	return accountInfo, nil
}
