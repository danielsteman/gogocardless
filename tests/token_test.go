package tests

import (
	"testing"

	"github.com/danielsteman/gogocardless/gocardless"
)

func TestCreateToken(t *testing.T) {
	_, err := gocardless.GetOrRefreshToken()
	if err != nil {
		t.Errorf("error getting token")
	}
}
