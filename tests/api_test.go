package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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

	// Test case: Code not provided
	t.Run("CodeNotProvided", func(t *testing.T) {
		body := map[string]string{
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

		expectedError := "code not provided"
		if response["error"] != expectedError {
			t.Errorf("Expected error %q, but got %q", expectedError, response["error"])
		}
	})

	// Test case: Not working code
	t.Run("NotWorkingCode", func(t *testing.T) {
		body := map[string]string{
			"code":     "print('Hello, World!)",
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

		recorder := httptest.NewRecorder()

		srv := server.NewServer()

		srv.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
		}

		var response map[string]string
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		if response["error"] == "" {
			t.Error("Expected an error message, but got none")
		}
	})

	// Test case: Authentication failure
	// Test case: Authentication failure
	t.Run("AuthenticationFailure", func(t *testing.T) {
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

		recorder := httptest.NewRecorder()

		srv := server.NewServer()

		// Set the API key check enabled and provide an incorrect API key
		os.Setenv("API_KEY_CHECK_ENABLED", "true")
		os.Setenv("API_KEY", "your-api-key")
		defer os.Unsetenv("API_KEY_CHECK_ENABLED")
		defer os.Unsetenv("API_KEY")

		req.Header.Set("X-Api-Key", "incorrect-api-key")

		srv.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code %d, but got %d", http.StatusUnauthorized, recorder.Code)
		}

		var response map[string]string
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		expectedError := "unauthorized"
		if response["error"] != expectedError {
			t.Errorf("Expected error %q, but got %q", expectedError, response["error"])
		}
	})

	// Test case: Successful authentication
	t.Run("SuccessfulAuthentication", func(t *testing.T) {
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

		recorder := httptest.NewRecorder()

		srv := server.NewServer()

		// Set the API key check enabled and provide the correct API key
		os.Setenv("API_KEY_CHECK_ENABLED", "true")
		os.Setenv("API_KEY", "valid-api-key")
		defer os.Unsetenv("API_KEY_CHECK_ENABLED")
		defer os.Unsetenv("API_KEY")

		req.Header.Set("X-Api-Key", "valid-api-key")

		srv.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
		}

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

	// Test case: API key check disabled
	t.Run("ApiKeyCheckDisabled", func(t *testing.T) {
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

		recorder := httptest.NewRecorder()

		srv := server.NewServer()

		// Set the API key check disabled
		os.Setenv("API_KEY_CHECK_ENABLED", "false")
		defer os.Unsetenv("API_KEY_CHECK_ENABLED")

		srv.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
		}

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
	// Test case: Invalid JSON payload
	t.Run("InvalidJSONPayload", func(t *testing.T) {
		requestBody := []byte(`{"code": "print('Hello, World!')", "language": "python",}`) // Invalid JSON payload

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

		expectedError := "invalid request body"
		if response["error"] != expectedError {
			t.Errorf("Expected error %q, but got %q", expectedError, response["error"])
		}
	})
	// Test case: Wrong HTTP method
	t.Run("WrongHTTPMethod", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/execute", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		srv := server.NewServer()

		srv.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status code %d, but got %d", http.StatusMethodNotAllowed, recorder.Code)
		}

		var response map[string]string
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		expectedError := "method not allowed"
		if response["error"] != expectedError {
			t.Errorf("Expected error %q, but got %q", expectedError, response["error"])
		}
	})
}
