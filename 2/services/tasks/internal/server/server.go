package server

import (
	"fmt"
	"net/http"
	"time"

	"app/services/tasks/internal/auth"
	"app/services/tasks/internal/handlers"
	"app/services/tasks/internal/middleware"
	"app/services/tasks/internal/store"
	shared_middleware "app/shared/middleware"
)

func NewServer(port, authGRPCAddr string) *http.Server {
	taskStore := store.NewMemoryStore()
	grpcClient, err := auth.NewGRPCAuthClient(authGRPCAddr, 3*time.Second)
	if err != nil {
		panic(err)
	}
	taskHandler := handlers.NewTaskHandler(taskStore)

	mux := http.NewServeMux()
	mux.Handle("POST /tasks", middleware.Auth(grpcClient)(http.HandlerFunc(taskHandler.Create)))
	mux.Handle("GET /tasks", middleware.Auth(grpcClient)(http.HandlerFunc(taskHandler.List)))
	mux.Handle("GET /tasks/{id}", middleware.Auth(grpcClient)(http.HandlerFunc(taskHandler.Get)))
	mux.Handle("PATCH /tasks/{id}", middleware.Auth(grpcClient)(http.HandlerFunc(taskHandler.Update)))
	mux.Handle("DELETE /tasks/{id}", middleware.Auth(grpcClient)(http.HandlerFunc(taskHandler.Delete)))

	handler := shared_middleware.RequestID(shared_middleware.Logging(mux))

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}
}
