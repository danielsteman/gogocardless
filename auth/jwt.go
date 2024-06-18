package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielsteman/gogocardless/config"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(config.Config.JWTSecret)

type JWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
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
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				http.Error(w, "Invalid signing method", http.StatusUnauthorized)
			}
			return jwtSecret, nil
		})

		println(claims.Email)

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(r *http.Request) *JWTClaims {
	if claims, ok := r.Context().Value("user").(*JWTClaims); ok {
		return claims
	}
	return nil
}
