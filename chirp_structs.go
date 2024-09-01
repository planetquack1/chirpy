package main

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

type Login struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

type UserWithoutPassword struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// type chirpValid struct {
// 	Valid bool `json:"valid"`
// }

// type chirpCleanedBody struct {
// 	CleanedBody string `json:"cleaned_body"`
// }
