package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"app/services/auth/internal/service"
	"app/shared/api/auth"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthServer(authService *service.AuthService) *AuthServer {
	return &AuthServer{authService: authService}
}

func (s *AuthServer) Verify(ctx context.Context, req *auth.VerifyRequest) (*auth.VerifyResponse, error) {
	valid, subject := s.authService.Verify(req.Token)
	if !valid {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	return &auth.VerifyResponse{
		Valid:   true,
		Subject: subject,
	}, nil
}
