package tests

import (
	"log"
	"os"
	"testing"

	"github.com/danielsteman/gogocardless/db"
	"github.com/danielsteman/gogocardless/gocardless"
)

func TestMain(m *testing.M) {
	db, err := db.GetDB(
		db.DBConfig{
			DBName: "gogocardless-test",
			Port:   5431,
		},
	)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	err = db.AutoMigrate(&gocardless.Token{})
	if err != nil {
		panic("failed to migrate database")
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestCreateToken(t *testing.T) {
	_, err := gocardless.GetOrRefreshToken()
	if err != nil {
		t.Errorf("error getting token")
	}
}
