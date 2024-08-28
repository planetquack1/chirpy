package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	// an error will be thrown if the JSON is invalid or has the wrong types
	// any missing fields will simply have their values in the struct set to their zero value

	cError := chirpError{Error: msg}

	dat, err := json.Marshal(cError)
	if err != nil {
		log.Printf("Error marshalling error JSON: %s", err)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}

var profaneWordList = []string{"kerfuffle", "sharbert", "fornax"}

func respondWithJSON(w http.ResponseWriter, code int, msg string) {

	// Clean the body
	words := strings.Fields(msg) // replace the words in the slice
	for i, word := range words {
		for _, profanity := range profaneWordList {
			if strings.EqualFold(word, profanity) {
				words[i] = "****"
				break
			}
		}
	}

	// Join the words back into a cleaned message
	cleanedMsg := strings.Join(words, " ")

	// Send the cleaned body
	cCleanedBody := chirpCleanedBody{CleanedBody: cleanedMsg}

	dat, err := json.Marshal(cCleanedBody)
	if err != nil {
		log.Printf("Error marshalling error JSON: %s", err)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}
