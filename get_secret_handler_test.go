package main

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetSecretHandler(t *testing.T) {
	var writer *httptest.ResponseRecorder
	writer = httptest.NewRecorder()

	reqBody := "{\"plain_text\": \"super secret\"}"
	request, _ := http.NewRequest("POST", "/", strings.NewReader(reqBody))
	mux.ServeHTTP(writer, request)

	hasher := md5.New()
	hasher.Write([]byte("super secret"))
	hv := hex.EncodeToString(hasher.Sum(nil))

	request = httptest.NewRequest("GET", "/"+hv, strings.NewReader(""))
	writer = httptest.NewRecorder()
	mux.ServeHTTP(writer, request)

	expectedResponseBody := "{\"data\":\"super secret\"}"

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	var resp = string(writer.Body.Bytes())
	if resp != expectedResponseBody {
		t.Errorf("Response body is %s", resp)
	}

	writer = httptest.NewRecorder()
	mux.ServeHTTP(writer, request)

	expectedResponseBody = "{\"data\":\"\"}"

	if writer.Code != 404 {
		t.Errorf("Response code is %v", writer.Code)
	}
	resp = string(writer.Body.Bytes())
	if resp != expectedResponseBody {
		t.Errorf("Response body is %s", resp)
	}
}
