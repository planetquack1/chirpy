package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var expires_in_seconds = 2

func (cfg *Config) createToken() (string, error) {

	// Get current time
	currentTime := time.Now()

	// Calculate expiration date
	expiresIn := time.Duration(expires_in_seconds) * time.Second

	// Convert times
	issuedAt := jwt.NewNumericDate(currentTime)
	expiresAt := jwt.NewNumericDate(currentTime.Add(expiresIn))

	registeredClaims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)

	// Sign the token using the secret key
	tokenString, err := token.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
