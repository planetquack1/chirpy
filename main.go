package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	checkForDebug()

	// Create new database
	db, err := NewDB("database.json")
	if err != nil {
		fmt.Println("Error: could not connect to database")
		os.Exit(1)
	}

	// Create new config struct
	cfg := &Config{}

	// API struct as database + config
	api := API{
		Config: cfg,
		DB:     db,
	}

	// Initialize ServeMux
	mux := http.NewServeMux()

	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")

	cfg.jwtSecret = jwtSecret

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", cfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /assets", assetHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", cfg.resetHandler)

	mux.HandleFunc("GET /api/chirps/{chirpID}", db.getChirpByID)
	mux.HandleFunc("POST /api/chirps", db.postChirp)

	mux.HandleFunc("POST /api/users", db.postUser)
	mux.HandleFunc("PUT /api/users", api.updateUser)

	mux.HandleFunc("POST /api/login", api.postLogin)

	// mux.HandleFunc("POST /api/refresh", api.postRefresh)

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	srv.ListenAndServe()

}
