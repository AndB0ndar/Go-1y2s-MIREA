package server

import (
	"app/services/tasks/internal/auth"
	"app/services/tasks/internal/handlers"
	"app/services/tasks/internal/middleware"
	"app/services/tasks/internal/store"
	shared_middleware "app/shared/middleware"
	"fmt"
	"net/http"
	"time"
)

func NewServer(port, authBaseURL string) *http.Server {
	taskStore := store.NewMemoryStore()

	authClient := auth.NewClient(authBaseURL, 3*time.Second)

	taskHandler := handlers.NewTaskHandler(taskStore)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", middleware.Auth(authClient)(taskHandler.Create))
	mux.HandleFunc("GET /tasks", middleware.Auth(authClient)(taskHandler.List))
	mux.HandleFunc("GET /tasks/{id}", middleware.Auth(authClient)(taskHandler.Get))
	mux.HandleFunc("PATCH /tasks/{id}", middleware.Auth(authClient)(taskHandler.Update))
	mux.HandleFunc("DELETE /tasks/{id}", middleware.Auth(authClient)(taskHandler.Delete))

	handler := shared_middleware.RequestID(shared_middleware.Logging(mux))

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}
}
