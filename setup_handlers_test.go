package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/borkod/secrets-app/handlers"
)

var mux *http.ServeMux

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	os.Exit(code)
}

func setUp() {
	mux = http.NewServeMux()
	mux.HandleFunc("/healthcheck", handlers.HealthCheckHandler)
	mux.HandleFunc("/", handlers.SecretHandler)
}

func TestHealthcheck(t *testing.T) {
	var writer *httptest.ResponseRecorder
	writer = httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/healthcheck", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	var resp = string(writer.Body.Bytes())
	if resp != http.StatusText(http.StatusOK) {
		t.Errorf("Response body is not 'OK'. Received %s", resp)
	}
}

func TestHandleGet(t *testing.T) {
	var writer *httptest.ResponseRecorder
	writer = httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	mux.ServeHTTP(writer, request)
	if writer.Code != 404 {
		t.Errorf("Response code is %v", writer.Code)
	}
	var resp handlers.SecretData
	json.Unmarshal(writer.Body.Bytes(), &resp)

	if resp.Data != "" {
		t.Errorf("Cannot retrieve JSON post")
	}
}

func TestHandlePost(t *testing.T) {
	var writer *httptest.ResponseRecorder
	writer = httptest.NewRecorder()
	content := strings.NewReader("")
	request, _ := http.NewRequest("POST", "/", content)
	mux.ServeHTTP(writer, request)

	if writer.Code != 400 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestHandlePut(t *testing.T) {
	var writer *httptest.ResponseRecorder
	writer = httptest.NewRecorder()
	content := strings.NewReader("")
	request, _ := http.NewRequest("PUT", "/", content)
	mux.ServeHTTP(writer, request)

	if writer.Code != 405 {
		t.Errorf("Response code is %v", writer.Code)
	}
}
