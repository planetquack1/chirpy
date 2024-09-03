package main

import "time"

type Chirp struct {
	// these tags indicate how the keys in the JSON should be mapped to the struct fields
	// the struct fields must be exported (start with a capital letter) if you want them parsed
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type chirpBody struct {
	Body string `json:"body"`
}

type chirpError struct {
	Error string `json:"error"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenInfo struct {
	ExpiresAt time.Time `json:"expires_at"` // TODO: change to date format
	Email     string    `json:"email"`
}

type Token struct {
	Token string `json:"token"`
}

type Login struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
}

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Password     []byte `json:"password"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type UserWithoutPassword struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// type chirpValid struct {
// 	Valid bool `json:"valid"`
// }

// type chirpCleanedBody struct {
// 	CleanedBody string `json:"cleaned_body"`
// }
