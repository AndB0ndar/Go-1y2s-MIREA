package main

import (
	"database/sql"
	"fmt"
	"os"

	"app/services/tasks/internal/repository"
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

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	repo := repository.NewPostgresRepo(db)

	srv := server.NewServer(port, authGRPCAddr, repo, log)
	log.WithField("port", port).Info("Tasks service starting")
	if err := srv.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("server failed")
	}
}
