package server

import (
	"encoding/json"
	"net/http"
	"os"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("API_KEY_CHECK_ENABLED") == "true" {
			apiKey := r.Header.Get("X-Api-Key")
			expectedApiKey := os.Getenv("API_KEY")
			if expectedApiKey == "" {
				errorResponse(w, "API key not set", http.StatusInternalServerError)
				return
			}
			if apiKey != expectedApiKey {
				errorResponse(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func errorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := map[string]string{
		"error": message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
