package gocardless

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/danielsteman/gogocardless/config"
	"github.com/danielsteman/gogocardless/db"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
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

type BaseRequisition struct {
	ID           string         `gorm:"type:uuid;primaryKey" json:"id"`
	Redirect     string         `gorm:"type:varchar(255)" json:"redirect"`
	Status       string         `gorm:"type:varchar(50)" json:"status"`
	Agreement    string         `gorm:"type:varchar(36)" json:"agreement"`
	Accounts     pq.StringArray `gorm:"type:text[]" json:"accounts"`
	Reference    string         `gorm:"type:varchar(100);unique" json:"reference"`
	UserLanguage string         `gorm:"type:varchar(10)" json:"user_language"`
	Link         string         `gorm:"type:varchar(255)" json:"link"`
}

type Requisition struct {
	BaseRequisition
}

type DBRequisition struct {
	gorm.Model
	BaseRequisition
	Email string `gorm:"type:varchar(255)" json:"email"`
}

func (DBRequisition) TableName() string {
	return "requisitions"
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

func GetEndUserRequisitionLink(institutionID string, email string) (Requisition, error) {
	userAgreement, err := GetEndUserAgreement(institutionID)
	if err != nil {
		return Requisition{}, fmt.Errorf("failed to get user agreement: %w", err)
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
		return Requisition{}, fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return Requisition{}, fmt.Errorf("error creating request: %w", err)
	}

	token, err := GetOrRefreshToken()
	if err != nil {
		return Requisition{}, fmt.Errorf("failed to get token: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Access)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Requisition{}, fmt.Errorf("failed to get redirect info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return Requisition{}, fmt.Errorf("failed to get redirect info: status code %d, response: %s", resp.StatusCode, string(body))
	}

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return Requisition{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var requisition Requisition
	err = json.Unmarshal(jsonData, &requisition)
	if err != nil {
		return Requisition{}, fmt.Errorf("failed to unmarshal redirect info: %w", err)
	}

	dbRequisition := DBRequisition{
		BaseRequisition: requisition.BaseRequisition,
		Email:           email,
	}

	if _, err := dbCreateRequisition(dbRequisition); err != nil {
		return Requisition{}, fmt.Errorf("error saving new requisition: %w", err)
	}

	return requisition, nil
}

func dbCreateRequisition(requisition DBRequisition) (string, error) {
	db, err := db.GetDB()
	if err != nil {
		return "", fmt.Errorf("error connecting to the database: %w", err)
	}

	if err := db.Create(&requisition).Error; err != nil {
		return "", fmt.Errorf("error creating requisition: %w", err)
	}

	return "Requisition created successfully", nil
}

func DBGetRequisition(ID string) (Requisition, error) {
	db, err := db.GetDB()
	if err != nil {
		return Requisition{}, fmt.Errorf("error connecting to the database: %w", err)
	}

	var requisition Requisition
	if err := db.First(&requisition, "id = ?", ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return Requisition{}, fmt.Errorf("requisition not found: %w", err)
		}
		return Requisition{}, fmt.Errorf("error retrieving requisition: %w", err)
	}

	return requisition, nil
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
