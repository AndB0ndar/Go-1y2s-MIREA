package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/sirupsen/logrus"

	"app/services/auth/internal/service"
	"app/shared/api/auth"
	"app/shared/middleware"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	authService *service.AuthService
	logger      *logrus.Logger
}

func NewAuthServer(authService *service.AuthService, logger *logrus.Logger) *AuthServer {
	return &AuthServer{authService: authService, logger: logger}
}

func (s *AuthServer) Verify(ctx context.Context, req *auth.VerifyRequest) (*auth.VerifyResponse, error) {
	var requestID string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get("x-request-id"); len(vals) > 0 {
			requestID = vals[0]
		}
	}
	ctx = context.WithValue(ctx, middleware.RequestIDKey, requestID)
	log := s.logger.WithField("request_id", requestID).WithField("component", "grpc_verify")
	log.Debug("received verify request")

	valid, subject := s.authService.Verify(req.Token)
	if !valid {
		log.WithField("valid", false).Info("token invalid")
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	log.WithField("valid", true).Info("token verified")
	return &auth.VerifyResponse{
		Valid:   true,
		Subject: subject,
	}, nil
}
