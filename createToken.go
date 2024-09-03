package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var token_expires_in_seconds = 3600
var refresh_token_expires_in_days = 60

func (cfg *Config) createToken(userID int) (string, error) {

	// Get current time
	currentTime := time.Now()

	// Calculate expiration date
	expiresIn := time.Duration(token_expires_in_seconds) * time.Second

	// Convert times
	issuedAt := jwt.NewNumericDate(currentTime)
	expiresAt := jwt.NewNumericDate(currentTime.Add(expiresIn))

	registeredClaims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		Subject:   fmt.Sprint(userID),
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

// func (cfg *Config) createRefreshTokenInfo(userID int) (RefreshTokenInfo, error) { // unfinished

// 	// Get current time
// 	currentTime := time.Now()

// 	// Calculate expiration duration
// 	expiresIn := time.Duration(token_expires_in_seconds) * time.Second

// 	// Convert times
// 	issuedAt := jwt.NewNumericDate(currentTime)
// 	expiresAt := jwt.NewNumericDate(currentTime.Add(expiresIn))

// 	registeredClaims := &jwt.RegisteredClaims{
// 		Issuer:    "chirpy",
// 		IssuedAt:  issuedAt,
// 		ExpiresAt: expiresAt,
// 		Subject:   fmt.Sprint(userID),
// 	}

// 	return RefreshTokenInfo{
// 		ExpiresAt: expiresIn,
// 	}, nil
// }
