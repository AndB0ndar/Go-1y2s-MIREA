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

	"app/services/auth/internal/grpc"
	"app/services/auth/internal/server"
	"app/services/auth/internal/service"
	pb "app/shared/api/auth"
)

func main() {
	// HTTP server
	port := os.Getenv("AUTH_PORT")
	if port == "" {
		port = "8081"
	}
	httpServer := server.NewServer(port)
	go func() {
		log.Printf("Auth service starting on port %s", port)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
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
			log.Fatalf("failed to listen: %v", err)
		}
		pb.RegisterAuthServiceServer(grpcServer, grpc.NewAuthServer(authService))
		log.Printf("Auth gRPC service starting on port %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
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
