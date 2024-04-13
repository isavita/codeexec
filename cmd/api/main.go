package main

import (
	"log"
	"net/http"

	"github.com/isavita/codeexec/cmd/api/server"
)

func main() {
	// Create a new API server
	srv := server.NewServer()

	// Start the HTTP server
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", srv))
}
