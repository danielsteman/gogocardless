package tests

import (
	"log"
	"os"
	"testing"

	"github.com/danielsteman/gogocardless/config"
	"github.com/danielsteman/gogocardless/db"
	"github.com/danielsteman/gogocardless/gocardless"
)

func TestMain(m *testing.M) {
	config.ResetAppConfig()
	config.LoadAppConfig("../.env.test")

	db, err := db.GetDB()

	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	err = db.AutoMigrate(
		&gocardless.Token{},
		&gocardless.DBRequisition{},
		&gocardless.AccountInfo{},
	)
	if err != nil {
		panic("failed to migrate database")
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestCreateToken(t *testing.T) {
	_, err := gocardless.GetOrRefreshToken()
	if err != nil {
		t.Errorf("error getting token: %v", err)
	}
}

func TestGetOrRefreshToken(t *testing.T) {
	firstToken, err := gocardless.GetOrRefreshToken()
	if err != nil {
		t.Errorf("error getting token first time: %v", err)
	}

	secondToken, err := gocardless.GetOrRefreshToken()
	if err != nil {
		t.Errorf("error getting token second time: %v", err)
	}

	if firstToken.Access != secondToken.Access {
		t.Errorf("expected %q and %q to be the same", firstToken.Access, secondToken.Access)
	}
}
