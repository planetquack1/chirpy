package main

import (
	"net/http"
)

func (cfg *Config) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r) // Pass the request to the next handler
	})
}
