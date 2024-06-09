package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/danielsteman/gogocardless/config"
	"github.com/danielsteman/gogocardless/gocardless"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
)

type userResource struct{}

func (rs userResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/redirect", userRedirectHandler)
	r.Get("/accounts", userAccountsHandler)

	return r
}

var jwtSecret = []byte(config.Config.JWTSecret)

type JWTClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext retrieves the user information from the request context
func GetUserFromContext(r *http.Request) *JWTClaims {
	if claims, ok := r.Context().Value("user").(*JWTClaims); ok {
		return claims
	}
	return nil
}

func userRedirectHandler(w http.ResponseWriter, r *http.Request) {
	institutionID := r.URL.Query().Get("institutionId")
	if institutionID == "" {
		http.Error(w, "institutionId query parameter is required", http.StatusBadRequest)
		return
	}

	user := GetUserFromContext(r)

	redirectInfo, err := gocardless.GetEndUserRequisitionLink(institutionID, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(redirectInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func userAccountsHandler(w http.ResponseWriter, r *http.Request) {
	agreementRef := r.URL.Query().Get("agreementRef")
	if agreementRef == "" {
		http.Error(w, "agreementRef query parameter is required", http.StatusBadRequest)
		return
	}

	accountInfo, err := gocardless.GetEndUserAccountInfo(agreementRef)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(accountInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
