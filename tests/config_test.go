package tests

import (
	"testing"

	"github.com/danielsteman/gogocardless/config"
)

func TestGetLocalDBURL(t *testing.T) {
	got := config.Config.DBURL
	want := "postgresql://admin:admin@localhost:5432/gogocardless-test"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
