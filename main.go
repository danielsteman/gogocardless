package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type RequestBody struct {
	SecretID string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secretID := os.Getenv("SECRET_ID")
	secretKey := os.Getenv("SECRET_KEY")

	client := &http.Client{}
	URL := "https://bankaccountdata.gocardless.com/api/v2/token/new/"

	requestBody := RequestBody {
		SecretID:   secretID,
		SecretKey:  secretKey,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal("Error marshaling request body:", err)
		return
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Fatal("Error creating request:", err)
		return
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response:", string(responseBody))
}