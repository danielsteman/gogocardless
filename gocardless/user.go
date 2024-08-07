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

type Requisition struct {
	ID           string         `gorm:"type:uuid;primaryKey" json:"id"`
	Redirect     string         `gorm:"type:varchar(255)" json:"redirect"`
	Status       string         `gorm:"type:varchar(50)" json:"status"`
	Agreement    string         `gorm:"type:varchar(36)" json:"agreement"`
	Accounts     pq.StringArray `gorm:"type:text[]" json:"accounts"`
	Reference    string         `gorm:"type:varchar(100);unique" json:"reference"`
	UserLanguage string         `gorm:"type:varchar(10)" json:"user_language"`
	Link         string         `gorm:"type:varchar(255)" json:"link"`
}

type DBRequisition struct {
	gorm.Model
	ID           string         `gorm:"type:uuid;primaryKey" json:"id"`
	Redirect     string         `gorm:"type:varchar(255)" json:"redirect"`
	Status       string         `gorm:"type:varchar(50)" json:"status"`
	Agreement    string         `gorm:"type:varchar(36)" json:"agreement"`
	Accounts     pq.StringArray `gorm:"type:text[]" json:"accounts"`
	Reference    string         `gorm:"type:varchar(100);unique" json:"reference"`
	UserLanguage string         `gorm:"type:varchar(10)" json:"user_language"`
	Link         string         `gorm:"type:varchar(255)" json:"link"`
	Email        string         `gorm:"type:varchar(255)" json:"email"`
}

func (DBRequisition) TableName() string {
	return "requisitions"
}

type AccountInfo struct {
	gorm.Model
	ID         string         `gorm:"type:uuid;primaryKey" json:"id"`
	Status     string         `gorm:"type:varchar(255)" json:"status"`
	Agreements string         `gorm:"type:varchar(255)" json:"agreements"`
	Accounts   pq.StringArray `gorm:"type:text[]" json:"accounts"`
	Reference  string         `gorm:"type:varchar(255)" json:"reference"`
}

type Account struct {
	ID   uint   `gorm:"primaryKey"`
	Iban string `json:"iban"`
}

type TransactionAmount struct {
	ID       uint   `gorm:"primaryKey"`
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type BookedTransaction struct {
	ID                                uint              `gorm:"primaryKey"`
	TransactionId                     string            `json:"transactionId"`
	DebtorName                        string            `json:"debtorName,omitempty"`
	DebtorAccountID                   uint              `json:"-"`
	DebtorAccount                     *Account          `json:"debtorAccount,omitempty"`
	TransactionAmountID               uint              `json:"-"`
	TransactionAmount                 TransactionAmount `json:"transactionAmount"`
	BookingDate                       string            `json:"bookingDate"`
	ValueDate                         string            `json:"valueDate"`
	RemittanceInformationUnstructured string            `json:"remittanceInformationUnstructured"`
	BankTransactionCode               string            `json:"bankTransactionCode,omitempty"`
	AccountInfoID                     string            `json:"-"` // Foreign key for AccountInfo
}

type PendingTransaction struct {
	ID                                uint              `gorm:"primaryKey"`
	TransactionAmountID               uint              `json:"-"`
	TransactionAmount                 TransactionAmount `json:"transactionAmount"`
	ValueDate                         string            `json:"valueDate"`
	RemittanceInformationUnstructured string            `json:"remittanceInformationUnstructured"`
	AccountInfoID                     string            `json:"-"` // Foreign key for AccountInfo
}

type Transactions struct {
	ID        uint                 `gorm:"primaryKey"`
	Booked    []BookedTransaction  `json:"booked" gorm:"foreignKey:AccountInfoID;references:ID"`
	Pending   []PendingTransaction `json:"pending" gorm:"foreignKey:AccountInfoID;references:ID"`
	AccountID string               `json:"-"` // Foreign key for AccountInfo
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
		ID:           requisition.ID,
		Redirect:     requisition.Redirect,
		Status:       requisition.Status,
		Agreement:    requisition.Agreement,
		Accounts:     requisition.Accounts,
		Reference:    requisition.Reference,
		UserLanguage: requisition.UserLanguage,
		Link:         requisition.Link,
		Email:        email,
	}

	if _, err := DBCreateRequisition(dbRequisition); err != nil {
		return Requisition{}, fmt.Errorf("error saving new requisition: %w", err)
	}

	return requisition, nil
}

func DBCreateRequisition(requisition DBRequisition) (string, error) {
	db, err := db.GetDB()
	if err != nil {
		return "", fmt.Errorf("error connecting to the database: %w", err)
	}

	if err := db.Create(&requisition).Error; err != nil {
		return "", fmt.Errorf("error creating requisition: %w", err)
	}

	return "Requisition created successfully", nil
}

func DBGetRequisition(value string, searchBy string) (DBRequisition, error) {
	db, err := db.GetDB()
	if err != nil {
		return DBRequisition{}, fmt.Errorf("error connecting to the database: %w", err)
	}

	var requisition DBRequisition
	var query string

	switch searchBy {
	case "id":
		query = "id = ?"
	case "email":
		query = "email = ?"
	default:
		return DBRequisition{}, fmt.Errorf("invalid search parameter: %s", searchBy)
	}

	if err := db.First(&requisition, query, value).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return DBRequisition{}, fmt.Errorf("requisition not found with %s: %w", searchBy, err)
		}
		return DBRequisition{}, fmt.Errorf("error retrieving requisition with %s: %w", searchBy, err)
	}

	return requisition, nil
}

func DBGetAccountInfo(agreementID string) (AccountInfo, error) {
	db, err := db.GetDB()
	if err != nil {
		return AccountInfo{}, fmt.Errorf("error connecting to the database: %w", err)
	}

	var accountInfo AccountInfo

	result := db.Where("id = ?", agreementID).First(&accountInfo)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return AccountInfo{}, nil
		}
		return AccountInfo{}, fmt.Errorf("error retrieving account information: %w", result.Error)
	}

	return accountInfo, nil
}

func DBCreateAccountInfo(accountInfo AccountInfo) (string, error) {
	db, err := db.GetDB()
	if err != nil {
		return "", fmt.Errorf("error connecting to the database: %w", err)
	}

	// Check if the account info with the given ID already exists
	var existingAccountInfo AccountInfo
	result := db.Where("id = ?", accountInfo.ID).First(&existingAccountInfo)
	if result.Error == nil {
		// Record with the same ID already exists
		return "", fmt.Errorf("account info with ID %s already exists", accountInfo.ID)
	} else if result.Error != gorm.ErrRecordNotFound {
		// An error occurred while querying the database
		return "", fmt.Errorf("error checking for existing account info: %w", result.Error)
	}

	// Create the new account info record
	result = db.Create(&accountInfo)
	if result.Error != nil {
		return "", fmt.Errorf("error creating account information: %w", result.Error)
	}

	return "Account info created successfully", nil
}

func GetEndUserAccountInfo(agreementID string, email string) (AccountInfo, error) {
	accountInfo, err := DBGetAccountInfo(agreementID)
	if err != nil {
		return AccountInfo{}, err
	}

	// Check if accountInfo already exists and its status is not "LN"
	if accountInfo.Status == "LN" {
		return accountInfo, nil
	}

	// If accountInfo is empty, proceed to fetch and create it
	if accountInfo.Agreements == "" {
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

		err = json.Unmarshal(jsonData, &accountInfo)
		if err != nil {
			return AccountInfo{}, fmt.Errorf("failed to unmarshal account info: %w", err)
		}

		_, err = DBCreateAccountInfo(accountInfo)
		if err != nil {
			return AccountInfo{}, fmt.Errorf("error saving new account info: %w", err)
		}
	}

	return accountInfo, nil
}

func GetEndUserTransactions(accountId string, email string) (Transactions, error) {
	// get agreements from `requisitions` that have status `LN`
	// that can be used to pull the latest data for a user
	url := fmt.Sprintf("https://bankaccountdata.gocardless.com/api/v2/accounts/%s/transactions/", accountId)

	db, err := db.GetDB()
	if err != nil {
		return Transactions{}, fmt.Errorf("error connecting to the database: %w", err)
	}
	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Transactions{}, fmt.Errorf("error creating request: %w", err)
	}

	// Get or refresh the token
	token, err := GetOrRefreshToken()
	if err != nil {
		return Transactions{}, fmt.Errorf("failed to get token: %w", err)
	}

	// Set headers
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Access)

	// Create an HTTP client and make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Transactions{}, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Transactions{}, fmt.Errorf("failed to get transactions: status code %d, response: %s", resp.StatusCode, string(body))
	}
	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return Transactions{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the response into the Transactions struct
	var transactions Transactions
	err = json.Unmarshal(jsonData, &transactions)
	if err != nil {
		return Transactions{}, fmt.Errorf("failed to unmarshal transactions: %w", err)
	}

	transactions.AccountID = accountId

	err = db.Transaction(func(tx *gorm.DB) error {
		// Delete existing transactions for the account
		if err := tx.Where("account_id = ?", accountId).Delete(&BookedTransaction{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing booked transactions: %w", err)
		}
		if err := tx.Where("account_id = ?", accountId).Delete(&PendingTransaction{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing pending transactions: %w", err)
		}

		// Create new booked transactions
		for _, booked := range transactions.Booked {
			booked.AccountInfoID = accountId
			if err := tx.Create(&booked).Error; err != nil {
				return fmt.Errorf("failed to create booked transaction: %w", err)
			}
		}

		// Create new pending transactions
		for _, pending := range transactions.Pending {
			pending.AccountInfoID = accountId
			if err := tx.Create(&pending).Error; err != nil {
				return fmt.Errorf("failed to create pending transaction: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return Transactions{}, fmt.Errorf("failed to store transactions in database: %w", err)
	}

	return transactions, nil
}
