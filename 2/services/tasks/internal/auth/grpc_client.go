package auth

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"app/shared/api/auth"
)

type GRPCAuthClient struct {
	client auth.AuthServiceClient
}

func NewGRPCAuthClient(addr string, timeout time.Duration) (*GRPCAuthClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(timeout))
	if err != nil {
		return nil, err
	}
	return &GRPCAuthClient{
		client: auth.NewAuthServiceClient(conn),
	}, nil
}

func (c *GRPCAuthClient) Verify(ctx context.Context, token string) (bool, error) {
	req := &auth.VerifyRequest{Token: token}
	resp, err := c.client.Verify(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return false, err
		}
		switch st.Code() {
		case codes.Unauthenticated:
			return false, nil
		case codes.DeadlineExceeded, codes.Unavailable:
			return false, ErrAuthUnavailable
		default:
			return false, err
		}
	}
	return resp.Valid, nil
}

var ErrAuthUnavailable = &AuthError{msg: "auth service unavailable"}

type AuthError struct{ msg string }

func (e *AuthError) Error() string { return e.msg }
