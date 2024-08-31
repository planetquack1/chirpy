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

type Email struct {
	Email string `json:"email"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// type chirpValid struct {
// 	Valid bool `json:"valid"`
// }

// type chirpCleanedBody struct {
// 	CleanedBody string `json:"cleaned_body"`
// }
