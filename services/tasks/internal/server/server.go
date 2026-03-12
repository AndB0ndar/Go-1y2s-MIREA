package server

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"app/services/tasks/internal/auth"
	"app/services/tasks/internal/handlers"
	"app/services/tasks/internal/metrics"
	"app/services/tasks/internal/middleware"
	"app/services/tasks/internal/repository"
	shared_middleware "app/shared/middleware"
)

func NewServer(
	port, authGRPCAddr, instanceID string,
	repo repository.TaskRepository,
	log *logrus.Logger,
	rabbitConn *amqp.Connection,
	queueName string,
) *http.Server {
	grpcClient, err := auth.NewGRPCAuthClient(authGRPCAddr, 3*time.Second, log)
	if err != nil {
		log.WithError(err).Fatal("failed to create auth client")
	}

	taskHandler := handlers.NewTaskHandler(repo, log, rabbitConn, queueName)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.Health())
	mux.HandleFunc("GET /metrics", metrics.MetricsHandler().ServeHTTP)

	authMiddleware := middleware.Auth(grpcClient, log)
	csrfMiddleware := middleware.CSRFProtection(log)

	mux.Handle("GET /tasks", authMiddleware(http.HandlerFunc(taskHandler.List)))
	mux.Handle("GET /tasks/{id}", authMiddleware(http.HandlerFunc(taskHandler.Get)))

	mux.Handle("GET /tasks/search", authMiddleware(http.HandlerFunc(taskHandler.Search)))
	mux.Handle("GET /tasks/searchvulnerable", authMiddleware(http.HandlerFunc(taskHandler.SearchVulnerable)))

	postHandler := authMiddleware(csrfMiddleware(http.HandlerFunc(taskHandler.Create)))
	mux.Handle("POST /tasks", postHandler)
	patchHandler := authMiddleware(csrfMiddleware(http.HandlerFunc(taskHandler.Update)))
	mux.Handle("PATCH /tasks/{id}", patchHandler)
	deleteHandler := authMiddleware(csrfMiddleware(http.HandlerFunc(taskHandler.Delete)))
	mux.Handle("DELETE /tasks/{id}", deleteHandler)

	handler := shared_middleware.RequestID(mux)
	handler = shared_middleware.SecurityHeaders(handler)
	handler = shared_middleware.AccessLog(log)(handler)
	handler = metrics.MetricsMiddleware(handler)
	handler = middleware.InstanceID(instanceID)(handler)

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}
}
