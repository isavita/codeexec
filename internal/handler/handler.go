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
