package server

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"app/services/auth/internal/handlers"
	"app/shared/middleware"
)

func NewServer(port string, log *logrus.Logger) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/login", handlers.Login)
	mux.HandleFunc("GET /auth/verify", handlers.Verify)

	handler := middleware.RequestID(mux)
	handler = middleware.AccessLog(log)(handler)

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}
}
