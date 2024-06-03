package gocardless

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/danielsteman/gogocardless/config"
	"github.com/google/uuid"
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

type RequisitionPayload struct {
	Redirect      string `json:"redirect"`
	InstitutionID string `json:"institution_id"`
	Reference     string `json:"reference"`
	Agreement     string `json:"agreement"`
	UserLanguage  string `json:"user_language"`
}

type Status struct {
	Short       string `json:"short"`
	Long        string `json:"long"`
	Description string `json:"description"`
}

type RedirectInfo struct {
	ID           string   `json:"id"`
	Redirect     string   `json:"redirect"`
	Status       Status   `json:"status"`
	Agreement    string   `json:"agreement"`
	Accounts     []string `json:"accounts"`
	Reference    string   `json:"reference"`
	UserLanguage string   `json:"user_language"`
	Link         string   `json:"link"`
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

func GetEndUserRequisitionLink(institutionID string) (RedirectInfo, error) {
	userAgreement, err := GetEndUserAgreement(institutionID)
	newReference := uuid.New().String()
	url := "https://bankaccountdata.gocardless.com/api/v2/requisitions/"
	if err != nil {
		log.Fatal("failed to get user agreement:", err)
	}
	payload := RequisitionPayload{
		Redirect:      config.Config.RedirectURL,
		InstitutionID: userAgreement.InstitutionID,
		Reference:     newReference,
		Agreement:     userAgreement.ID,
		UserLanguage:  "EN",
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
		log.Fatal("failed to get redirect info:", err)
	}

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return RedirectInfo{}, fmt.Errorf("failed to get redirect info: status code %d, response: %s", resp.StatusCode, string(body))
	}

	defer resp.Body.Close()
	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("failed to read response body:", err)
	}

	var redirectInfo RedirectInfo
	err = json.Unmarshal([]byte(jsonData), &redirectInfo)
	if err != nil {
		log.Fatal("failed to get redirect info:", err)
	}

	return redirectInfo, nil
}
