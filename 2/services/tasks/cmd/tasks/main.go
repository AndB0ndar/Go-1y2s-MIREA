package main

import (
	"app/services/tasks/internal/server"
	"log"
	"os"
)

func main() {
	port := os.Getenv("TASKS_PORT")
	if port == "" {
		port = "8082"
	}
	authGRPCAddr := os.Getenv("AUTH_GRPC_ADDR")
	if authGRPCAddr == "" {
		authGRPCAddr = "localhost:50051"
	}
	srv := server.NewServer(port, authGRPCAddr)
	log.Printf("Tasks service starting on port %s, auth gRPC addr: %s", port, authGRPCAddr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
