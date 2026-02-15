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
	authBaseURL := os.Getenv("AUTH_BASE_URL")
	if authBaseURL == "" {
		authBaseURL = "http://localhost:8081"
	}
	srv := server.NewServer(port, authBaseURL)
	log.Printf("Tasks service starting on port %s, auth URL: %s", port, authBaseURL)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
