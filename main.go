package main

import (
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	apiCfg := &apiConfig{}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mux.HandleFunc("/assets", assetHandler)
	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/reset", apiCfg.resetHandler)

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	srv.ListenAndServe()

}
