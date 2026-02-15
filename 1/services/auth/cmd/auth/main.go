package main

import (
	"log"
	"os"
	"app/services/auth/internal/server"
)

func main() {
	port := os.Getenv("AUTH_PORT")
	if port == "" {
		port = "8081"
	}
	srv := server.NewServer(port)
	log.Printf("Auth service starting on port %s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
