package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *Config) getUserIDFromToken(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	fmt.Println("token string is " + tokenString)

	claims, err := cfg.parseToken(tokenString)
	if err != nil {
		fmt.Println("error parsing token")
		return "", err
	}

	return claims.Subject, nil
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
