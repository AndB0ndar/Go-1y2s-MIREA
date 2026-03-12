package auth

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "app/shared/api/auth"
)

type Client struct {
	conn    *grpc.ClientConn
	client  pb.AuthServiceClient
	timeout time.Duration
}

func NewClient(addr string, timeout time.Duration) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(timeout))
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:    conn,
		client:  pb.NewAuthServiceClient(conn),
		timeout: timeout,
	}, nil
}

func (c *Client) Verify(ctx context.Context, token string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	req := &pb.VerifyRequest{Token: token}
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

func (c *Client) Close() error {
	return c.conn.Close()
}

var ErrAuthUnavailable = &AuthError{msg: "auth service unavailable"}

type AuthError struct{ msg string }

func (e *AuthError) Error() string { return e.msg }
