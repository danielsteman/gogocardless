package tests

import (
	"testing"

	"github.com/danielsteman/gogocardless/config"
)

func TestGetSecretID(t *testing.T) {
	got := config.Config.SecretID
	want := "example-secret-id"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestGetLocalDBURL(t *testing.T) {
	got := config.Config.DBURL
	want := "localhost:420"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
