package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func assetHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Length", "35672")
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func (cfg *Config) metricsHandler(w http.ResponseWriter, r *http.Request) {
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

func (cfg *Config) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func (db *DB) getChirps(w http.ResponseWriter, r *http.Request) {

	authorIDstr := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")

	database, err := db.loadDB()
	if err != nil {
		respondWithError(w, 500, "Error loading database")
		return
	}

	// If authorID is not provided, get all chirps
	chirps := make(map[int]Chirp)
	if authorIDstr == "" {
		chirps = database.Chirps
	} else {

		// Convert authorID to integer using strconv.Atoi
		authorID, err := strconv.Atoi(authorIDstr)
		if err != nil {
			respondWithError(w, 400, "Invalid chirpID") // 400 error?
			return
		}

		// Search through all chirps, find where author IDs match
		for _, chirp := range database.Chirps {
			if chirp.AuthorID == authorID {
				chirps[chirp.Id] = chirp
			}
		}
	}

	reversed := false
	if sort == "desc" {
		reversed = true
	} else if sort == "asc" || sort == "" {
		reversed = false
	} else {
		respondWithError(w, 400, "Invalid sort method") // 400 error?
		return
	}

	chirpsList := chirpsToList(chirps, reversed)

	dat, err := json.Marshal(chirpsList)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

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

func (api *API) postChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	cBody := chirpBody{}

	// Check if cannot decode, send ERROR
	err := decoder.Decode(&cBody)
	if err != nil {
		respondWithError(w, 500, "Cannot decode body into struct")
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

	// Extract JWT token from the Authorization header
	userID, err := api.Config.getUserIDFromToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token, or cannot parse ID as int")
		return
	}

	// Load the database
	database, err := api.DB.loadDB()
	if err != nil {
		respondWithError(w, 400, "Could not load database")
		return
	}

	// Create the cleaned Chirp struct
	cCleanedChirp := Chirp{
		Id:       len(database.Chirps) + 1,
		Body:     cleanedMsg,
		AuthorID: userID,
	}

	// Marshal the cleaned chirp struct
	dat, err := json.Marshal(cCleanedChirp)
	if err != nil {
		log.Printf("Error marshalling cleaned chirp: %s", err)
		return
	}

	// Add post to local database
	database.Chirps[len(database.Chirps)+1] = cCleanedChirp

	// Write to the original database
	if err := api.DB.writeDB(database); err != nil {
		respondWithError(w, 500, "Error saving database")
		return
	}

	// Write to HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)

}

func (api *API) deleteChirp(w http.ResponseWriter, r *http.Request) {

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

	// Extract JWT token from the Authorization header
	userID, err := api.Config.getUserIDFromToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token, or cannot parse ID as int")
		return
	}

	// Load the database
	database, err := api.DB.loadDB()
	if err != nil {
		respondWithError(w, 400, "Could not load database")
		return
	}

	// Check if chirp exists
	if _, exists := database.Chirps[chirpID]; !exists {
		respondWithError(w, http.StatusInternalServerError, "Chirp does not exist")
		return
	}

	// Check if user is the author of the chirp
	if database.Chirps[chirpID].AuthorID != userID {
		respondWithError(w, 403, "Cannot delete another user's chirp")
		return
	}

	// Clear chirp
	database.Chirps[chirpID] = Chirp{
		Id:       database.Chirps[chirpID].Id, // Can keep the same
		Body:     "",
		AuthorID: 0,
	}

	// Write to the original database
	if err := api.DB.writeDB(database); err != nil {
		respondWithError(w, 500, "Error saving database")
		return
	}

	// Write to HTTP response
	w.WriteHeader(http.StatusNoContent)
}

func (db *DB) postUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	login := Login{}

	// Check if cannot decode, send ERROR
	err := decoder.Decode(&login)
	if err != nil {
		respondWithError(w, 500, "Cannot decode body into struct")
		return
	}

	// Load current database
	database, err := db.loadDB()
	if err != nil {
		respondWithError(w, 500, "Error loading database")
		return
	}

	// Encrypt password
	cost := bcrypt.DefaultCost

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(login.Password), cost)
	if err != nil {
		respondWithError(w, 500, "Error generating password hash")
		return
	}

	// Create new User struct
	user := User{
		ID:           len(database.Users) + 1,
		Email:        login.Email,
		Password:     encryptedPassword,
		RefreshToken: "",
		IsChirpyRed:  false,
	}

	// Check if user exists in database
	if _, exists := database.Users[login.Email]; exists {
		respondWithError(w, 500, "User with same email exists")
		return
	}
	// Add user to local database
	database.Users[login.Email] = user
	database.UsersByID = append(database.UsersByID, user.Email)

	// Write to the original database
	if err := db.writeDB(database); err != nil {
		respondWithError(w, 500, "Error saving database")
		return
	}

	// Create new User Without Password struct
	uWithoutPassword := UserWithoutPassword{
		ID:    len(database.Users),
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
	w.WriteHeader(201)
	w.Write(dat)

}

func (api *API) updateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	updatedLogin := Login{}

	// Check if cannot decode, send ERROR
	err := decoder.Decode(&updatedLogin)
	if err != nil {
		respondWithError(w, 500, "Cannot decode body into struct")
		return
	}

	// VALIDATE USER

	// Extract JWT token from the Authorization header
	userID, err := api.Config.getUserIDFromToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token, or cannot parse ID as int")
		return
	}

	// UPDATE USER

	// Load current database
	database, err := api.DB.loadDB()
	if err != nil {
		respondWithError(w, 500, "Error loading database")
		return
	}

	// Check if user exists in database
	userEmail := database.UsersByID[userID-1] // Get email as key
	user, exists := database.Users[userEmail] // Use key to access user struct
	if !exists {
		respondWithError(w, http.StatusUnauthorized, "User does not exist")
		return
	}

	// Check if overwriting another user's information
	if contains(database.UsersByID, updatedLogin.Email) && updatedLogin.Email != userEmail {
		respondWithError(w, http.StatusUnauthorized, "User with same email exists")
		return
	}

	// Encrypt password
	cost := bcrypt.DefaultCost

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedLogin.Password), cost)
	if err != nil {
		respondWithError(w, 500, "Error generating password hash")
		return
	}

	// VOID THE PREVIOUS TOKEN

	previousToken, exists := database.RefreshTokens[user.RefreshToken]
	if !exists {
		respondWithError(w, 500, "Cannot find refresh token")
		return
	}
	database.RefreshTokens[user.RefreshToken] = RefreshTokenInfo{
		ExpiresAt: previousToken.ExpiresAt, // Keep the same, no need to change
		Email:     "",
	}

	// Create an updated User struct
	updatedUser := User{
		ID:           userID,
		Email:        updatedLogin.Email,
		Password:     encryptedPassword,
		RefreshToken: user.RefreshToken, // Keep same refresh token
		IsChirpyRed:  user.IsChirpyRed,
	}

	// TODO: Delete user with old email as its key. For now, set all fields to empty

	// Get the old email (key)
	userEmailBeforeChange := database.UsersByID[userID-1]
	// Set fields to empty
	database.Users[userEmailBeforeChange] = User{
		ID:       0,
		Email:    "",
		Password: []byte(""),
	}

	// Update both user lists
	database.UsersByID[userID-1] = updatedLogin.Email
	database.Users[updatedLogin.Email] = updatedUser

	// Update email of refresh token, if exists
	refreshToken := user.RefreshToken
	fmt.Println("user token: " + refreshToken)
	if existingTokenInfo, exists := database.RefreshTokens[refreshToken]; exists {
		database.RefreshTokens[refreshToken] = RefreshTokenInfo{
			ExpiresAt: existingTokenInfo.ExpiresAt, // Keep the same ExpiresAt value
			Email:     updatedLogin.Email,          // Update the email
		}
	}

	// Write to the original database
	if err := api.DB.writeDB(database); err != nil {
		respondWithError(w, 500, "Error saving database")
		return
	}

	// WRITE TO HTTP RESPONSE

	// Create an updated User Without Password struct
	uWithoutPassword := UserWithoutPassword{
		ID:    userID,
		Email: updatedLogin.Email,
	}

	// Marshal the UserWithoutPassword struct
	dat, err := json.Marshal(uWithoutPassword)
	if err != nil {
		log.Printf("Error marshalling user: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)

}

func (api *API) postLogin(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	login := Login{}

	// Check if cannot decode, send ERROR
	err := decoder.Decode(&login)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot decode request body")
		return
	}

	// Load current database
	database, err := api.DB.loadDB()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error loading database")
		return
	}

	// VALIDATE USER

	// Check if user exists in database
	user, exists := database.Users[login.Email]
	if !exists {
		respondWithError(w, http.StatusUnauthorized, "User does not exist")
		return
	}

	// See if password is correct
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(login.Password)); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Password is incorrect")
		return
	}

	// VOID THE PREVIOUS TOKEN

	previousToken, exists := database.RefreshTokens[user.RefreshToken]
	if exists {
		database.RefreshTokens[user.RefreshToken] = RefreshTokenInfo{
			ExpiresAt: previousToken.ExpiresAt, // Keep the same, no need to change
			Email:     "",
		}

	}

	// UPDATE USER WITH NEW REFRESH TOKEN

	// Generate a refresh token
	refreshTokenBytes := make([]byte, 32)
	_, err = rand.Read(refreshTokenBytes)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}
	refreshToken := hex.EncodeToString(refreshTokenBytes)

	database.Users[user.Email] = User{
		ID:           user.ID,
		Email:        user.Email,
		Password:     user.Password,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	}

	// ADD TOKEN TO DATABASE

	// Calculate expiration duration
	expiresIn := time.Duration(refresh_token_expires_in_days) * 24 * time.Hour

	database.RefreshTokens[refreshToken] = RefreshTokenInfo{
		ExpiresAt: time.Now().Add(expiresIn),
		Email:     user.Email,
	}

	// Write to the original database
	if err := api.DB.writeDB(database); err != nil {
		respondWithError(w, 500, "Error saving database")
		return
	}

	// RESPOND WITH USER WITHOUT PASSWORD

	// Create a token
	token, err := api.Config.createToken(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Password is incorrect")
		return
	}

	// Create UserWithoutPassword struct
	uWithoutPassword := UserWithoutPassword{
		ID:           user.ID,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	}

	// Marshal the UserWithoutPassword struct
	dat, err := json.Marshal(uWithoutPassword)
	if err != nil {
		log.Printf("Error marshalling user: %s", err)
		return
	}

	// Write to HTTP response
	w.Header().Set("Authorization", token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}

func (api *API) postRefresh(w http.ResponseWriter, r *http.Request) {

	// Extract refresh token from the Authorization header
	refreshToken := getTokenFromHeader(r, "Bearer")

	// Load current database
	database, err := api.DB.loadDB()
	if err != nil {
		respondWithError(w, 500, "Error loading database")
		return
	}

	// Check if token exists in database
	refreshTokenInfo, exists := database.RefreshTokens[refreshToken]
	if !exists {
		respondWithError(w, 500, "Refresh token not found")
		return
	}

	// Get the user
	user, exists := database.Users[refreshTokenInfo.Email]
	if !exists {
		respondWithError(w, 500, "Cannot find owner of refresh token")
		return
	}

	// Token is valid, so create a new access token
	token, err := api.Config.createToken(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Password is incorrect")
		return
	}

	// WRITE TO HTTP RESPONSE

	// Create new Token struct
	tokenResponse := Token{
		Token: token,
	}

	// Marshal the Token struct
	dat, err := json.Marshal(tokenResponse)
	if err != nil {
		log.Printf("Error marshalling user: %s", err)
		return
	}

	w.Header().Set("Authorization", token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)

}

func (api *API) postRevoke(w http.ResponseWriter, r *http.Request) {

	// Extract refresh token from the Authorization header
	refreshToken := getTokenFromHeader(r, "Bearer")

	// Load current database
	database, err := api.DB.loadDB()
	if err != nil {
		respondWithError(w, 500, "Error loading database")
		return
	}

	// Check if token exists in database
	refreshTokenInfo, exists := database.RefreshTokens[refreshToken]
	if !exists {
		respondWithError(w, 500, "Refresh token not found")
		return
	}

	// Get the user
	user, exists := database.Users[refreshTokenInfo.Email]
	if !exists {
		respondWithError(w, http.StatusUnauthorized, "Cannot find owner of refresh token")
		return
	}

	// VOID THE PREVIOUS TOKEN

	// Refresh token exists, so delete it, TODO: remove entirely
	database.RefreshTokens[refreshToken] = RefreshTokenInfo{
		ExpiresAt: refreshTokenInfo.ExpiresAt, // keep this the same
		Email:     "",
	}

	// Remove old refresh token in Users list
	database.Users[refreshTokenInfo.Email] = User{
		ID:           user.ID,
		Email:        user.Email,
		Password:     user.Password,
		RefreshToken: "", // Remove refresh token
		IsChirpyRed:  user.IsChirpyRed,
	}

	// Write to the original database
	if err := api.DB.writeDB(database); err != nil {
		respondWithError(w, 500, "Error saving database")
		return
	}

	// WRITE TO HTTP RESPONSE

	w.WriteHeader(http.StatusNoContent)

}

func (api *API) postPolka(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	polka := Polka{}

	// Check if cannot decode, send ERROR
	err := decoder.Decode(&polka)
	if err != nil {
		respondWithError(w, 500, "Cannot decode body into struct")
		return
	}

	// Extract Polka API key from the Authorization header
	polkaAPIKey := getTokenFromHeader(r, "ApiKey")

	// Check if API key matches
	if polkaAPIKey != api.Config.polkaSecret {
		respondWithError(w, 401, "Unauthorized purchase")
		return
	}

	// Check if event is not user.upgraded, so do not edit database
	if polka.Event != "user.upgraded" {
		return
	}

	// Get user that upgraded
	userID := polka.Data.UserID

	// Load current database
	database, err := api.DB.loadDB()
	if err != nil {
		respondWithError(w, 500, "Error loading database")
		return
	}

	// Get user email by ID
	userEmail := database.UsersByID[userID-1]

	// Get user by email
	user, exists := database.Users[userEmail]
	if !exists {
		respondWithError(w, 404, "Could not find user by email")
		return
	}

	// Update user
	database.Users[userEmail] = User{
		ID:           user.ID,
		Email:        user.Email,
		Password:     user.Password,
		RefreshToken: user.RefreshToken,
		IsChirpyRed:  true,
	}

	// Write to the original database
	if err := api.DB.writeDB(database); err != nil {
		respondWithError(w, 500, "Error saving database")
		return
	}

	// Write to HTTP response
	w.WriteHeader(http.StatusNoContent)

}
