package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/borkod/secrets-app/fileHandler"
)

var secretsValues = secretsStore{mu: sync.Mutex{}, secrets: make(map[string]string), dataFileName: ""}

//
func InitiateSecretHandler(s string) error {
	fn, err := fileHandler.CheckCreateFile(s)
	if err != nil {
		return err
	}
	secretsValues.dataFileName = fn
	return nil
}

// SecretHandler processes both GET and POST requests
func SecretHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		secretHandlerPost(w, r)
	case http.MethodGet:
		secretHandlerGet(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// secretHandlerPost processes POST requests
func secretHandlerPost(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength == 0 {
		w.WriteHeader(400)
		fmt.Fprintln(w, "No secret provided")
		return
	}

	// parse request input to get value of "plain_text"
	sv, err := parseInput(r)
	if err != nil || len(sv) == 0 {
		w.WriteHeader(400)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}

	// store secret value in the map data structure
	secretId, err := storeSecret(sv)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	// send response back to client
	w.Header().Set("Content-Type", "application/json")
	resp := &SecretID{
		Id: secretId,
	}
	json, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "Error formatting response")
		return
	}
	w.Write(json)
}

// secretHandlerGet processes GET requests
func secretHandlerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := strings.TrimPrefix(r.URL.Path, "/")
	s, err := secretsValues.GetDeleteSecret(id)
	if s == "" {
		w.WriteHeader(404)
	}

	resp := &SecretData{
		Data: s,
	}
	json, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "Error formatting response")
		return
	}
	w.Write(json)

}

// parse request input to get value of "plain_text"
func parseInput(r *http.Request) (string, error) {
	var c PlainText
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	err := json.Unmarshal(body, &c)
	if err != nil {
		return "", err
	}
	return c.PlainText, nil
}

// updates map data structure to store secret value
func storeSecret(s string) (string, error) {
	return secretsValues.AddSecret(s)
}
