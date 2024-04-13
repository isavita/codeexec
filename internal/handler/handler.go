package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/isavita/codeexec/internal/executor"
)

type CodeExecutionHandler struct {
	executor *executor.DockerExecutor
}

func NewCodeExecutionHandler() *CodeExecutionHandler {
	exec, err := executor.NewDockerExecutor()
	if err != nil {
		panic(err)
	}
	return &CodeExecutionHandler{executor: exec}
}

func (h *CodeExecutionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body map[string]string
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	code := body["code"]
	language := body["language"]

	if language == "" {
		errorResponse(w, "language not specified", http.StatusBadRequest)
		return
	}

	if !isLanguageSupported(language) {
		errorResponse(w, "unsupported language: "+language, http.StatusBadRequest)
		return
	}

	output, err := h.executor.Execute(code, language, 5*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"output": output,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func isLanguageSupported(language string) bool {
	supportedLanguages := []string{"python", "javascript"}
	for _, lang := range supportedLanguages {
		if lang == language {
			return true
		}
	}
	return false
}

func errorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := map[string]string{
		"error": message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
