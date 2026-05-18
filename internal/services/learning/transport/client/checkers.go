package client

import (
	"context"
	"time"

	authpb "online-learning-platform/pkg/gen/auth"
	coursepb "online-learning-platform/pkg/gen/courses"

	"google.golang.org/grpc"
)

type AuthUserChecker struct {
	conn   *grpc.ClientConn
	client authpb.AuthServiceClient
}

func NewAuthUserChecker(addr string) (*AuthUserChecker, error) {
	conn, err := dial(addr)
	if err != nil {
		return nil, err
	}
	return &AuthUserChecker{conn: conn, client: authpb.NewAuthServiceClient(conn)}, nil
}

func (c *AuthUserChecker) Exists(ctx context.Context, userID string) (bool, error) {
	_, err := c.client.GetUserProfile(ctx, &authpb.GetUserProfileRequest{UserId: userID})
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (c *AuthUserChecker) Close() error { return c.conn.Close() }

type CourseChecker struct {
	conn   *grpc.ClientConn
	client coursepb.CourseServiceClient
}

func NewCourseChecker(addr string) (*CourseChecker, error) {
	conn, err := dial(addr)
	if err != nil {
		return nil, err
	}
	return &CourseChecker{conn: conn, client: coursepb.NewCourseServiceClient(conn)}, nil
}

func (c *CourseChecker) Exists(ctx context.Context, courseID string) (bool, error) {
	_, err := c.client.GetCourse(ctx, &coursepb.GetCourseRequest{Id: courseID})
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (c *CourseChecker) Close() error { return c.conn.Close() }

func dial(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
}
