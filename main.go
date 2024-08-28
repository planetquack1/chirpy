package main

import (
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	apiCfg := &apiConfig{}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /assets", assetHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.resetHandler)
	// mux.HandleFunc("GET /api/validate_chirp", apiCfg.validateChirpHandler)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirpHandler)

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	srv.ListenAndServe()

}
