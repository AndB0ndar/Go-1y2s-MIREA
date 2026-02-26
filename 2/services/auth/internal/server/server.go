package server

import (
	"fmt"
	"net/http"

	"app/services/auth/internal/handlers"
	"app/shared/middleware"
)

func NewServer(port string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/login", handlers.Login)
	mux.HandleFunc("GET /auth/verify", handlers.Verify)

	handler := middleware.RequestID(middleware.Logging(mux))

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}
}
