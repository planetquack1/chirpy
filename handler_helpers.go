package main

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *Config) getUserIDFromToken(r *http.Request) (int, error) {

	tokenString := getTokenFromHeader(r, "Bearer")

	fmt.Println("token string is " + tokenString)

	claims, err := cfg.parseToken(tokenString)
	if err != nil {
		fmt.Println("error parsing token")
		return 0, err
	}

	// Convert userID to int
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (cfg *Config) parseToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			fmt.Println("unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.jwtSecret), nil
	})

	if err != nil {
		fmt.Println("error")
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		fmt.Printf("got claims. Claims: ")
		fmt.Println(claims.Subject)
		return claims, nil
	} else {
		fmt.Println("invalid token")
		return nil, errors.New("invalid token")
	}
}

func getTokenFromHeader(r *http.Request, prefix string) string {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, (prefix + " ")) {
		// Return the token without the "Bearer " prefix
		return strings.TrimPrefix(authHeader, (prefix + " "))
	}
	// Return the header as is if "Bearer " is not found
	return authHeader
}

func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}

func chirpsToList(chirps map[int]Chirp, reversed bool) []Chirp {

	chirpsSlice := []Chirp{}

	for _, chirp := range chirps {
		chirpsSlice = append(chirpsSlice, chirp)
	}

	// Sort chirps by ID
	sort.Slice(chirpsSlice, func(i, j int) bool {
		return chirpsSlice[i].Id < chirpsSlice[j].Id
	})

	if reversed {
		for i, j := 0, len(chirpsSlice)-1; i < j; i, j = i+1, j-1 {
			chirpsSlice[i], chirpsSlice[j] = chirpsSlice[j], chirpsSlice[i]
		}
	}

	return chirpsSlice
}
