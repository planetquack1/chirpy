package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	checkForDebug()

	db, err := NewDB("database.json")
	if err != nil {
		fmt.Println("Error: could not connect to database")
		os.Exit(1)
	}

	mux := http.NewServeMux()
	apiCfg := &apiConfig{}

	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")

	apiCfg.jwtSecret = jwtSecret

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /assets", assetHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.resetHandler)

	mux.HandleFunc("GET /api/chirps/{chirpID}", db.getChirpByID)
	mux.HandleFunc("POST /api/chirps", db.postChirp)

	mux.HandleFunc("POST /api/users", db.postUser)

	mux.HandleFunc("POST /api/login", db.postLogin)

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	srv.ListenAndServe()

}
