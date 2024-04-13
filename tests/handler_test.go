package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isavita/codeexec/internal/handler"
)

func TestCodeExecutionHandler(t *testing.T) {
	// Create a new HTTP request with the code and language in the request body
	body := map[string]string{
		"code":     "print('Hello, Python!')",
		"language": "python",
	}
	requestBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/execute", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a new HTTP response recorder
	recorder := httptest.NewRecorder()

	// Create a new code execution handler
	h := handler.NewCodeExecutionHandler()

	// Serve the HTTP request to the handler
	h.ServeHTTP(recorder, req)

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

	expectedOutput := "Hello, Python!"
	if response["output"] != expectedOutput {
		t.Errorf("Expected output %q, but got %q", expectedOutput, response["output"])
	}
}
