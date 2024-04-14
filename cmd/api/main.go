package main

import (
	"log"
	"net/http"
	"os"

	"github.com/isavita/codeexec/cmd/api/server"
)

func main() {
	// Get the port from the environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a new API server
	srv := server.NewServer()

	// Start the HTTP server
	log.Printf("Server listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, srv))
}
