package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"

	"app/services/tasks/internal/cache"
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

	instanceID := os.Getenv("INSTANCE_ID")
	if instanceID == "" {
		instanceID = "unknown"
	}

	// DB
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
	var repo repository.TaskRepository
	repo = repository.NewPostgresRepo(db)

	// Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	cacheTTL, _ := strconv.Atoi(os.Getenv("CACHE_TTL_SECONDS"))
	if cacheTTL == 0 {
		cacheTTL = 120
	}
	cacheJitter, _ := strconv.Atoi(os.Getenv("CACHE_TTL_JITTER_SECONDS"))
	if cacheJitter == 0 {
		cacheJitter = 30
	}

	redisClient, err := cache.NewRedisClient(redisAddr, redisPassword, 0, log)
	if err != nil {
		log.WithError(err).Warn("Redis unavailable, continuing without cache")
	} else {
		repo = repository.NewCachedTaskRepository(
			repo, redisClient, log, cacheTTL, cacheJitter,
		)
	}

	srv := server.NewServer(port, authGRPCAddr, instanceID, repo, log)
	log.WithField("port", port).Info("Tasks service starting")
	if err := srv.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("server failed")
	}
}
