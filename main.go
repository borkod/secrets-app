package main

import (
	"net/http"

	"github.com/borkod/secrets-app/handlers"
)

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/healthcheck/", handlers.HealthCheckHandler)
	http.HandleFunc("/", handlers.SecretHandler)

	server.ListenAndServe()
}
