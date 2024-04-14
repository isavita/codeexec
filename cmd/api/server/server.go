package server

import (
	"net/http"

	"github.com/isavita/codeexec/internal/handler"
)

func NewServer() http.Handler {
	codeExecutionHandler := handler.NewCodeExecutionHandler()

	mux := http.NewServeMux()
	mux.Handle("/api/execute", AuthMiddleware(codeExecutionHandler))

	return mux
}
