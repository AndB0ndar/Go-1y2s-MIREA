package server

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"app/services/auth/internal/handlers"
	"app/services/auth/internal/service"
	"app/shared/middleware"
)

func NewServer(
	port string, authService *service.AuthService, log *logrus.Logger,
) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/login", handlers.Login(authService, log))
	mux.HandleFunc("GET /auth/verify", handlers.Verify(authService, log))

	mux.HandleFunc("GET /health", handlers.Health())

	// middleware: RequestID -> SecurityHeaders -> Logging -> ...
	handler := middleware.RequestID(mux)
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.AccessLog(log)(handler)

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}
}
