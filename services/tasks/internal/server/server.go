package server

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"app/services/tasks/internal/auth"
	"app/services/tasks/internal/handlers"
	"app/services/tasks/internal/metrics"
	"app/services/tasks/internal/middleware"
	"app/services/tasks/internal/repository"
	shared_middleware "app/shared/middleware"
)

func NewServer(
	port, authGRPCAddr string, repo repository.TaskRepository, log *logrus.Logger,
) *http.Server {
	grpcClient, err := auth.NewGRPCAuthClient(authGRPCAddr, 3*time.Second, log)
	if err != nil {
		log.WithError(err).Fatal("failed to create auth client")
	}

	taskHandler := handlers.NewTaskHandler(repo, log)

	mux := http.NewServeMux()

	mux.Handle("GET /tasks", middleware.Auth(grpcClient, log)(http.HandlerFunc(taskHandler.List)))
	mux.Handle("GET /tasks/{id}", middleware.Auth(grpcClient, log)(http.HandlerFunc(taskHandler.Get)))

	mux.Handle("GET /tasks/search", middleware.Auth(grpcClient, log)(http.HandlerFunc(taskHandler.Search)))
	mux.Handle("GET /tasks/searchvulnerable", middleware.Auth(grpcClient, log)(http.HandlerFunc(taskHandler.SearchVulnerable)))

	postHandler := middleware.Auth(grpcClient, log)(middleware.CSRFProtection(log)(http.HandlerFunc(taskHandler.Create)))
	mux.Handle("POST /tasks", postHandler)
	patchHandler := middleware.Auth(grpcClient, log)(middleware.CSRFProtection(log)(http.HandlerFunc(taskHandler.Update)))
	mux.Handle("PATCH /tasks/{id}", patchHandler)
	deleteHandler := middleware.Auth(grpcClient, log)(middleware.CSRFProtection(log)(http.HandlerFunc(taskHandler.Delete)))
	mux.Handle("DELETE /tasks/{id}", deleteHandler)

	mux.HandleFunc("GET /metrics", metrics.MetricsHandler().ServeHTTP)

	mux.HandleFunc("GET /health", handlers.HealthCheck(log))

	handler := shared_middleware.RequestID(mux)
	handler = shared_middleware.SecurityHeaders(handler)
	handler = shared_middleware.AccessLog(log)(handler)
	handler = metrics.MetricsMiddleware(handler)

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}
}
