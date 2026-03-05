package auth

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/sirupsen/logrus"

	"app/shared/api/auth"
	"app/shared/middleware"
)

type GRPCAuthClient struct {
	client auth.AuthServiceClient
	logger *logrus.Logger
}

func NewGRPCAuthClient(addr string, timeout time.Duration, log *logrus.Logger) (*GRPCAuthClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(timeout))
	if err != nil {
		return nil, err
	}
	return &GRPCAuthClient{
		client: auth.NewAuthServiceClient(conn),
		logger: log,
	}, nil
}

func (c *GRPCAuthClient) Verify(ctx context.Context, token string) (bool, error) {
	requestID := middleware.GetRequestID(ctx)
	c.logger.WithFields(logrus.Fields{
		"request_id": requestID,
		"component":  "auth_client",
	}).Debug("calling auth.Verify via gRPC")
	md := metadata.New(map[string]string{"x-request-id": requestID})
	outCtx := metadata.NewOutgoingContext(ctx, md)

	req := &auth.VerifyRequest{Token: token}
	resp, err := c.client.Verify(outCtx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.logger.WithError(err).WithField("request_id", requestID).Error("unexpected error from auth")
			return false, err
		}
		switch st.Code() {
		case codes.Unauthenticated:
			c.logger.WithField("request_id", requestID).Info("auth returned unauthenticated")
			return false, nil
		case codes.DeadlineExceeded, codes.Unavailable:
			c.logger.WithField("request_id", requestID).Warn("auth unavailable or deadline exceeded")
			return false, ErrAuthUnavailable
		default:
			c.logger.WithField("request_id", requestID).WithError(err).Error("auth returned error")
			return false, err
		}
	}
	c.logger.WithField("request_id", requestID).Debug("auth verify successful")
	return resp.Valid, nil
}

var ErrAuthUnavailable = &AuthError{msg: "auth service unavailable"}

type AuthError struct{ msg string }

func (e *AuthError) Error() string { return e.msg }
