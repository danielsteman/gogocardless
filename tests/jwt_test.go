package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/danielsteman/gogocardless/auth"
	"github.com/danielsteman/gogocardless/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestVerifyToken(t *testing.T) {
	createTestToken := func(email string, secret []byte) string {
		claims := &auth.JWTClaims{
			Email: email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(secret)
		return tokenString
	}

	tests := []struct {
		name          string
		authHeader    string
		expectedCode  int
		expectedEmail string
	}{
		{
			name:         "Missing Authorization header",
			authHeader:   "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Invalid token format",
			authHeader:   "InvalidFormat",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Invalid token",
			authHeader:   "Bearer invalidtoken",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:          "Valid token",
			authHeader:    "Bearer " + createTestToken("test@example.com", []byte(config.Config.JWTSecret)),
			expectedCode:  http.StatusOK,
			expectedEmail: "test@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request with the specified Authorization header
			req, err := http.NewRequest("GET", "/", nil)
			assert.NoError(t, err)
			req.Header.Set("Authorization", tt.authHeader)

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Create a test handler to pass to the middleware
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				claims := auth.GetUserFromContext(r)
				if claims != nil {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(claims.Email))
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			})

			// Wrap the test handler with the VerifyToken middleware
			handler := auth.VerifyToken(testHandler)

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedCode, rr.Code)

			// Check the email if the status is OK
			if tt.expectedCode == http.StatusOK {
				assert.Equal(t, tt.expectedEmail, rr.Body.String())
			}
		})
	}
}
