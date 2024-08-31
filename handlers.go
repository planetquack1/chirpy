package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func assetHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Length", "35672")
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	htmlResponse := fmt.Sprintf(`
    <html>
    <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited %d times!</p>
    </body>
    </html>
    `, cfg.fileserverHits)
	w.Write([]byte(htmlResponse))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

// func (db *DB) getAllChirps(w http.ResponseWriter, r *http.Request) {

// 	database, err := db.loadDB()
// 	if err != nil {
// 		respondWithError(w, 500, "Error loading database")
// 		return
// 	}

// 	dat, err := json.Marshal(database)
// 	if err != nil {
// 		log.Printf("Error marshalling JSON: %s", err)
// 		w.WriteHeader(500)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(200)
// 	w.Write(dat)
// }

func (db *DB) getChirpByID(w http.ResponseWriter, r *http.Request) {

	// Extract chirpID from the path using PathValue
	chirpIDStr := r.PathValue("chirpID")
	if chirpIDStr == "" {
		respondWithError(w, 400, "chirpID not provided") // 400 error?
		return
	}

	// Convert chirpID to integer using strconv.Atoi
	chirpID, err := strconv.Atoi(chirpIDStr)
	if err != nil {
		respondWithError(w, 400, "Invalid chirpID") // 400 error?
		return
	}

	// Load current database
	database, err := db.loadDB()
	if err != nil {
		respondWithError(w, 500, "Error loading database")
		return
	}

	// Look up the chirp in the database
	chirp, exists := database.Chirps[chirpID-1]
	if !exists {
		respondWithError(w, 500, "Chirp not found") // 500 error?
		return
	}

	// Marhsal the chirp
	dat, err := json.Marshal(chirp)
	if err != nil {
		log.Printf("Error marshalling error JSON: %s", err)
		return
	}

	// Return the chirp as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)

}

func (db *DB) postChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	cBody := chirpBody{}

	// Check if cannot decode, send ERROR
	err := decoder.Decode(&cBody)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	// Check if body is at most 140 characters
	if len(cBody.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	// Clean the body
	words := strings.Fields(cBody.Body) // replace the words in the slice
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

	// Load the database
	database, err := db.loadDB()
	if err != nil {
		respondWithError(w, 400, "Could not load database")
		return
	}

	// Create the cleaned Chirp struct
	cCleanedChirp := Chirp{
		Id:   len(database.Chirps) + 1,
		Body: cleanedMsg,
	}

	// Marshal the cleaned chirp struct
	dat, err := json.Marshal(cCleanedChirp)
	if err != nil {
		log.Printf("Error marshalling cleaned chirp: %s", err)
		return
	}

	// Add post to local database
	database.Chirps[len(database.Chirps)] = cCleanedChirp

	// Write to the original database
	if err := db.writeDB(database); err != nil {
		respondWithError(w, 500, "Error saving database")
		return
	}

	// Write to HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}

func (db *DB) postUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	login := Login{}

	// Check if cannot decode, send ERROR
	err := decoder.Decode(&login)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	// Load current database
	database, err := db.loadDB()
	if err != nil {
		respondWithError(w, 500, "Error loading database")
		return
	}

	// Create new User struct
	user := User{
		ID:       len(database.Users) + 1,
		Email:    login.Email,
		Password: login.Password,
	}

	// Check if user exists in database
	if _, exists := database.Users[login.Email]; exists {
		respondWithError(w, 500, "User with same email exists")
		return
	}
	// Add user to local database
	database.Users[login.Email] = user

	// Write to the original database
	if err := db.writeDB(database); err != nil {
		respondWithError(w, 500, "Error saving database")
		return
	}

	// Create new User Without Password struct
	uWithoutPassword := UserWithoutPassword{
		ID:    len(database.Users) + 1,
		Email: login.Email,
	}

	// Marshal the UserWithoutPassword struct
	dat, err := json.Marshal(uWithoutPassword)
	if err != nil {
		log.Printf("Error marshalling user: %s", err)
		return
	}

	// Write to HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}
