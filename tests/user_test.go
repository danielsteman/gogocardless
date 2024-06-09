package tests

import (
	"log"
	"testing"

	"github.com/danielsteman/gogocardless/gocardless"
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
	requisition, err := gocardless.DBGetRequisition(endUserAgreement.ID)
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
	accountInfo, err := gocardless.GetEndUserAccountInfo(agreementID)
	if err != nil {
		log.Fatalf("Error getting account info: %v", err)
	}
	if len(accountInfo.Accounts) == 0 {
		log.Fatalf("Did not get the expected number of accounts: %v", err)
	}
}
