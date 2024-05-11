package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Credentials struct {
	SecretID string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}

type Token struct {
	Access         string `json:"access"`
	AccessExpires  int    `json:"access_expires"`
	Refresh        string `json:"refresh"`
	RefreshExpires int    `json:"refresh_expires"`
    }


func NewHTTPClient(timeout time.Duration) *http.Client {
    return &http.Client{Timeout: timeout}
}

func GetCredentials() (Credentials, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secretID := os.Getenv("SECRET_ID")
	secretKey := os.Getenv("SECRET_KEY")
	if secretID == "" {
		return Credentials{}, fmt.Errorf("SECRET_ID environment variable is empty")
	}
	if secretKey == "" {
		return Credentials{}, fmt.Errorf("SECRET_KEY environment variable is empty")
	}

	return Credentials{
		SecretID: secretID,
		SecretKey: secretKey,
	}, nil
}

func GetToken(client *http.Client, credentials Credentials) (*Token, error) {
	URL := "https://bankaccountdata.gocardless.com/api/v2/token/new/"
	requestBodyBytes, err := json.Marshal(credentials)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var token Token
	if err := json.Unmarshal(responseBody, &token); err != nil {
		return nil, fmt.Errorf("error decoding token: %v", err)
	}

	return &token, nil
}

func GoCardLessManager() {
	credentials, err := GetCredentials()
	if err != nil {
		log.Fatal("Couldn't get credentials:", err)
		return
	}

	client := &http.Client{}
	token, err := GetToken(client, credentials)
	if err != nil {
		log.Fatal("Couldn't get token:", err)
		return
	}
	fmt.Println("Token:", string(token.Access))
}


func main() {
	GoCardLessManager()
}