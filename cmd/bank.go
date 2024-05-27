package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/danielsteman/gogocardless/config"
	"github.com/go-chi/chi/v5"
)

type bankResource struct{}

func (rs bankResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/list", ListBanks) // GET /users - read a list of users
	// r.Post("/", rs.Create) // POST /users - create a new user and persist it
	// r.Put("/", rs.Delete)

	// r.Route("/{id}", func(r chi.Router) {
	// 	// r.Use(rs.TodoCtx) // lets have a users map, and lets actually load/manipulate
	// 	r.Get("/", rs.Get)       // GET /users/{id} - read a single user by :id
	// 	r.Put("/", rs.Update)    // PUT /users/{id} - update a single user by :id
	// 	r.Delete("/", rs.Delete) // DELETE /users/{id} - delete a single user by :id
	// })

	return r
}

func ListBanks(w http.ResponseWriter, r *http.Request) {

}

type Credentials struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}

type Token struct {
	Access         string `json:"access"`
	AccessExpires  int    `json:"access_expires"`
	Refresh        string `json:"refresh"`
	RefreshExpires int    `json:"refresh_expires"`
}

func GetToken() (*Token, error) {
	config, _ := config.LoadAppConfig(".env")
	url := "https://bankaccountdata.gocardless.com/api/v2/token/new/"
	credentials := Credentials{
		SecretID:  config.SecretID,
		SecretKey: config.SecretKey,
	}
	credentailsData, err := json.Marshal(credentials)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(credentailsData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	client := &http.Client{}
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

// func GetOrRefreshToken() (*Token, error) {

// }
