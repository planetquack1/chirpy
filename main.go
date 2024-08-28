package main

import (
	"net/http"
)

func main() {

	mu := http.NewServeMux()

	srv := http.Server{
		Addr:    ":8080",
		Handler: mu,
	}

	srv.ListenAndServe()

}
