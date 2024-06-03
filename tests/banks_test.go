package tests

import (
	"log"
	"testing"

	"github.com/danielsteman/gogocardless/gocardless"
)

func TestGetBankList(t *testing.T) {
	banks, err := gocardless.ListBanks()
	if err != nil {
		log.Fatalf("Error listing banks: %v", err)
	}
	if len(banks) == 0 {
		log.Fatalf("Did not get any banks: %v", err)
	}
}
