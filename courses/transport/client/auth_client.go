package client

import (
	"context"
	"time"

	authpb "online-learning-platform/pkg/gen/auth"

	"google.golang.org/grpc"
)

type AuthClient struct {
	conn   *grpc.ClientConn
	client authpb.AuthServiceClient
}

func NewAuthClient(addr string) (*AuthClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return &AuthClient{conn: conn, client: authpb.NewAuthServiceClient(conn)}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}

func (c *AuthClient) CanManageCourses(ctx context.Context, userID string) (bool, error) {
	resp, err := c.client.GetUserProfile(ctx, &authpb.GetUserProfileRequest{UserId: userID})
	if err != nil {
		return false, err
	}
	return resp.GetRole() == "instructor" || resp.GetRole() == "admin", nil
}
