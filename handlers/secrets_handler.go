package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var secretsValues = secrets{mu: sync.Mutex{}, secrets: make(map[string]string)}

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

func secretHandlerPost(w http.ResponseWriter, r *http.Request) {
	sv, err := parseInput(r)
	if err != nil {
		w.WriteHeader(501)
		fmt.Fprintln(w, "Error parsing input")
	}

	secretId, err := storeSecret(sv)
	if err != nil {
		w.WriteHeader(501)
		fmt.Fprintln(w, "Error processing secret")
	}

	w.Header().Set("Content-Type", "application/json")
	resp := &SecretID{
		Id: secretId,
	}
	json, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(501)
		fmt.Fprintln(w, "Error formatting response")
	}
	w.Write(json)
}

func secretHandlerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := strings.TrimPrefix(r.URL.Path, "/")
	s := ""
	_, ok := secretsValues.secrets[id]
	if ok {
		s = secretsValues.GetDeleteSecret(id)
	} else {
		w.WriteHeader(404)
	}

	resp := &SecretData{
		Data: s,
	}
	json, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(501)
		fmt.Fprintln(w, "Error formatting response")
	}
	w.Write(json)

}

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

func storeSecret(s string) (string, error) {
	return secretsValues.AddSecret(s), nil
}
