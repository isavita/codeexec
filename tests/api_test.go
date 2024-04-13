package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isavita/codeexec/cmd/api/server"
)

func TestCodeExecutionEndpoint(t *testing.T) {
	// Test case: Successful execution
	t.Run("SuccessfulExecution", func(t *testing.T) {
		// Create a new HTTP request with the code and language in the request body
		body := map[string]string{
			"code":     "print('Hello, World!')",
			"language": "python",
		}
		requestBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}

		req, err := http.NewRequest("POST", "/api/execute", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Create a new HTTP response recorder
		recorder := httptest.NewRecorder()

		// Create a new API server
		server := server.NewServer()

		// Serve the HTTP request to the API server
		server.ServeHTTP(recorder, req)

		// Check the response status code
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
		}

		// Check the response body
		var response map[string]string
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		expectedOutput := "Hello, World!"
		if response["output"] != expectedOutput {
			t.Errorf("Expected output %q, but got %q", expectedOutput, response["output"])
		}
	})
	// Test case: Unsupported language
	t.Run("UnsupportedLanguage", func(t *testing.T) {
		body := map[string]string{
			"code":     "println('Hello, World!')",
			"language": "unsupported",
		}
		requestBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}

		req, err := http.NewRequest("POST", "/api/execute", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		srv := server.NewServer()

		srv.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, recorder.Code)
		}

		var response map[string]string
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		expectedError := "unsupported language: unsupported"
		if response["error"] != expectedError {
			t.Errorf("Expected error %q, but got %q", expectedError, response["error"])
		}
	})

	// Test case: Language not specified
	t.Run("LanguageNotSpecified", func(t *testing.T) {
		body := map[string]string{
			"code": "println('Hello, World!')",
		}
		requestBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}

		req, err := http.NewRequest("POST", "/api/execute", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		srv := server.NewServer()

		srv.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, recorder.Code)
		}

		var response map[string]string
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		expectedError := "language not specified"
		if response["error"] != expectedError {
			t.Errorf("Expected error %q, but got %q", expectedError, response["error"])
		}
	})
}
