package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var secretsValues secrets

type PlainText struct {
	PlainText string `json:"plain_text"`
}

type SecretID struct {
	Id string `json:"id"`
}

type SecretData struct {
	Data string `json:"data"`
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func secretHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		secretHandlerPost(w, r)
	default:
		secretHandlerGet(w, r)
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

type secrets struct {
	// mu guards the secrets map.
	mu      sync.Mutex
	secrets map[string]string
}

func (s *secrets) AddSecret(v string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	hasher := md5.New()
	hasher.Write([]byte(v))
	m := hex.EncodeToString(hasher.Sum(nil))
	s.secrets[m] = v
	return m
}

func (s *secrets) GetDeleteSecret(secret string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	val := secretsValues.secrets[secret]
	delete(secretsValues.secrets, secret)
	return val
}

func main() {
	secretsValues.secrets = make(map[string]string)
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/healthcheck/", healthCheckHandler)
	http.HandleFunc("/", secretHandler)

	server.ListenAndServe()
}
