package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"app/services/graphql/internal/server"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	port := os.Getenv("GRAPHQL_PORT")
	if port == "" {
		port = "8090"
	}
	authGRPCAddr := os.Getenv("AUTH_GRPC_ADDR")
	if authGRPCAddr == "" {
		authGRPCAddr = "auth:50051"
	}

	// DB
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connStr := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to db")
	}
	if err := db.Ping(); err != nil {
		log.WithError(err).Fatal("failed to ping db")
	}

	srv := server.NewServer(port, authGRPCAddr, db, log)
	log.WithField("port", port).Info("GraphQL service starting")
	if err := srv.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("server failed")
	}
}
