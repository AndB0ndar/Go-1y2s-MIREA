package main

import (
	"os"

	"app/services/tasks/internal/server"
	"app/shared/logger"
)

func main() {
	log := logger.New("tasks")

	port := os.Getenv("TASKS_PORT")
	if port == "" {
		port = "8082"
	}
	authGRPCAddr := os.Getenv("AUTH_GRPC_ADDR")
	if authGRPCAddr == "" {
		authGRPCAddr = "localhost:50051"
	}
	srv := server.NewServer(port, authGRPCAddr, log)
	log.WithField("port", port).Info("Tasks service starting")
	if err := srv.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("server failed")
	}
}
