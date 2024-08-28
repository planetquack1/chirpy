package main

type chirpBody struct {
	Body string `json:"body"`
}

// type chirpValid struct {
// 	Valid bool `json:"valid"`
// }

type chirpError struct {
	Error string `json:"error"`
}

type chirpCleanedBody struct {
	CleanedBody string `json:"cleaned_body"`
}
