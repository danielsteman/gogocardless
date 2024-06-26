package tests

import (
	"log"
	"testing"

	"github.com/danielsteman/gogocardless/db"
	"github.com/danielsteman/gogocardless/gocardless"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetEndUserAgreement(t *testing.T) {
	want := "RABOBANK_RABONL2U"
	endUserAgreement, err := gocardless.GetEndUserAgreement(want)
	if err != nil {
		log.Fatalf("Error getting end user agreement: %v", err)
	}
	if endUserAgreement.InstitutionID != want {
		log.Fatalf("Did not get the expected institution: %v", err)
	}
}

func TestGetEndUserRequisitionLink(t *testing.T) {
	want := "CR"
	institutionID := "RABOBANK_RABONL2U"
	endUserAgreement, err := gocardless.GetEndUserRequisitionLink(institutionID, "test@test.com")
	if err != nil {
		log.Fatalf("Error getting redirect info: %v", err)
	}
	if endUserAgreement.Status != want {
		log.Fatalf("Did not get the expected status: %v", err)
	}
	requisition, err := gocardless.DBGetRequisition(endUserAgreement.ID, "id")
	if err != nil {
		log.Fatalf("Error getting requisition from database: %v", err)
	}
	if requisition.ID != endUserAgreement.ID {
		log.Fatalf("Did not find created requisition in database: %v", err)
	}
}

func TestGetEndUserAccountInfo(t *testing.T) {
	t.Skip("Skipping testing with potentially invalid agreementID")
	agreementID := "1006584c-d7a8-4cc4-988c-32af67bf1d02"
	accountInfo, err := gocardless.GetEndUserAccountInfo(agreementID, "test@test.com")
	if err != nil {
		log.Fatalf("Error getting account info: %v", err)
	}
	if len(accountInfo.Accounts) == 0 {
		log.Fatalf("Did not get the expected number of accounts: %v", err)
	}
}

func TestPutAccountInfo(t *testing.T) {
	requisition := gocardless.DBRequisition{
		ID:           uuid.New().String(),
		Redirect:     "test1234",
		Status:       "test1234",
		Agreement:    "test1234",
		Accounts:     []string{"test1", "test2"},
		Reference:    "test1234",
		UserLanguage: "test1234",
		Link:         "test1234",
		Email:        "",
	}
	_, err := gocardless.DBCreateRequisition(requisition)
	if err != nil {
		t.Fatalf("Failed to create requisition: %v", err)
	}

	// Perform the update
	err = gocardless.DBPutRequisition(requisition.Agreement, "Status", "LN")
	if err != nil {
		t.Fatalf("Failed to update requisition: %v", err)
	}

	// Verify the update
	db, err := db.GetDB()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	var updatedRequisition gocardless.DBRequisition
	result := db.Where("agreement = ?", requisition.Agreement).First(&updatedRequisition)
	if result.Error != nil {
		t.Fatalf("Failed to retrieve updated requisition: %v", result.Error)
	}

	assert.Equal(t, "LN", updatedRequisition.Status, "The Status field should be updated to 'LN'")

	db.Delete(&updatedRequisition)
}
