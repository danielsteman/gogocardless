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
