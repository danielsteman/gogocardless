package tests

import (
	"testing"

	"github.com/danielsteman/gogocardless/config"
)

func TestAdd(t *testing.T) {
	config, _ := config.LoadAppConfig(".env.test")
	got := config.SecretID
	want := "example-secret-id"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
