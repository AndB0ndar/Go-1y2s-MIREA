package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	google_grpc "google.golang.org/grpc"

	"app/shared/logger"

	"app/services/auth/internal/grpc"
	"app/services/auth/internal/server"
	"app/services/auth/internal/service"
	pb "app/shared/api/auth"
)

func main() {
	logger := logger.New("auth")

	// HTTP server
	httpPort := os.Getenv("AUTH_PORT")
	if httpPort == "" {
		httpPort = "8081"
	}
	httpServer := server.NewServer(httpPort, logger)
	go func() {
		logger.WithField("port", httpPort).Info("Auth HTTP service starting")
		if err := httpServer.ListenAndServe(); err != nil {
			logger.WithError(err).Fatal("HTTP server failed")
		}
	}()

	// gRPC server
	grpcPort := os.Getenv("AUTH_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}
	grpcServer := google_grpc.NewServer()
	authService := service.NewAuthService()
	go func() {
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			logger.WithError(err).Fatal("failed to listen")
		}
		pb.RegisterAuthServiceServer(grpcServer, grpc.NewAuthServer(authService, logger))
		logger.WithField("port", grpcPort).Info("Auth gRPC service starting")
		if err := grpcServer.Serve(lis); err != nil {
			logger.WithError(err).Fatal("failed to serve")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	grpcServer.GracefulStop()

	log.Println("Shutting down servers...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server shutdown failed: %v", err)
	}
	log.Println("Servers stopped")
}
