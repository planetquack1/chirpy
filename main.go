package main

import (
	"fmt"
	"net/http"
	"os"
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

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /assets", assetHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.resetHandler)

	mux.HandleFunc("GET /api/chirps/{chirpID}", db.getChirpByID)
	mux.HandleFunc("POST /api/chirps", db.postChirp)

	mux.HandleFunc("POST /api/users", db.postUser)

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	srv.ListenAndServe()

}
