package tests

import (
	"fmt"
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

	err = db.AutoMigrate(&gocardless.Token{})
	if err != nil {
		panic("failed to migrate database")
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestCreateToken(t *testing.T) {
	token, err := gocardless.GetOrRefreshToken()
	fmt.Println("Token:", string(token.Access))
	if err != nil {
		t.Errorf("error getting token")
	}
}
