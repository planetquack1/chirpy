package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// type chirp struct {
// 	// these tags indicate how the keys in the JSON should be mapped to the struct fields
// 	// the struct fields must be exported (start with a capital letter) if you want them parsed
// 	Body  string `json:"body"`
// 	Valid bool   `json:"valid"`
// 	Error string `json:"error"`
// }

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

// func (cfg *apiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {

// 	c := chirp{
// 		Body:  "HI",
// 		Valid: true,
// 		Error: "",
// 	}

// 	dat, err := json.Marshal(c)
// 	if err != nil {
// 		log.Printf("Error marshalling JSON: %s", err)
// 		c.Error = "Error marshalling JSON"
// 		w.WriteHeader(500)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(200)
// 	w.Write(dat)
// }

func (cfg *apiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

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

	respondWithJSON(w, 200, cBody.Body)

}
