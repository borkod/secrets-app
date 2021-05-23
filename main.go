package main

import (
	"log"
	"net/http"
	"os"

	"github.com/borkod/secrets-app/handlers"
)

func main() {

	path := os.Getenv("DATA_FILE_PATH")
	if len(path) == 0 {
		log.Fatalln("Please set DATA_FILE_PATH environment variable.")
	}

	if err := handlers.InitiateSecretHandler(path); err != nil {
		log.Fatalln(err.Error())
	}
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/healthcheck", handlers.HealthCheckHandler)
	http.HandleFunc("/", handlers.SecretHandler)

	server.ListenAndServe()
}
