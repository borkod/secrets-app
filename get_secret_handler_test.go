package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetSecretHandler(t *testing.T) {
	type testConfig struct {
		requestBody            string
		requestAction          string
		expectedHTTPStatusCode int
		expectedBody           string
	}

	testCases := []testConfig{
		testConfig{
			requestBody:            "{\"plain_text\": \"super secret\"}",
			requestAction:          "POST",
			expectedHTTPStatusCode: 200,
			expectedBody:           fmt.Sprintf("{\"id\":\"%s\"}", fmt.Sprintf("%x", md5.Sum([]byte("super secret")))),
		},
		testConfig{
			requestBody:            fmt.Sprintf("%x", md5.Sum([]byte("super secret"))),
			requestAction:          "GET",
			expectedHTTPStatusCode: 200,
			expectedBody:           "{\"data\":\"super secret\"}",
		},
		testConfig{
			requestBody:            fmt.Sprintf("%x", md5.Sum([]byte("super secret"))),
			requestAction:          "GET",
			expectedHTTPStatusCode: 404,
			expectedBody:           "{\"data\":\"\"}",
		},
	}

	for _, tc := range testCases {
		var writer *httptest.ResponseRecorder
		writer = httptest.NewRecorder()
		var request *http.Request
		if tc.requestAction == "GET" {
			request, _ = http.NewRequest("GET", "/"+tc.requestBody, strings.NewReader(""))
		} else {
			request, _ = http.NewRequest("POST", "/", strings.NewReader(tc.requestBody))
		}
		mux.ServeHTTP(writer, request)

		if writer.Code != tc.expectedHTTPStatusCode {
			t.Errorf("Response code is |%v|", writer.Code)
		}

		resp := string(writer.Body.Bytes())
		if resp != tc.expectedBody {
			t.Errorf("Response body is |%s|", resp)
		}
	}
}
