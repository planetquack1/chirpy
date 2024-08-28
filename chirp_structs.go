package main

type chirp struct {
	// these tags indicate how the keys in the JSON should be mapped to the struct fields
	// the struct fields must be exported (start with a capital letter) if you want them parsed
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type chirpBody struct {
	Body string `json:"body"`
}

// type chirpValid struct {
// 	Valid bool `json:"valid"`
// }

type chirpError struct {
	Error string `json:"error"`
}

// type chirpCleanedBody struct {
// 	CleanedBody string `json:"cleaned_body"`
// }
