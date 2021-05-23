package handlers

import (
	"net/http"
)

// Verifies that the application is ready to handle requests
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
