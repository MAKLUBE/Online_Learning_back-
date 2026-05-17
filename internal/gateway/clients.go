package gateway

import (
	"context"
	"os"
	"time"

	authpb "online-learning-platform/pkg/gen/auth"
	coursepb "online-learning-platform/pkg/gen/courses"
	learningpb "online-learning-platform/pkg/gen/learning"
	notificationpb "online-learning-platform/pkg/gen/notifications"

	"google.golang.org/grpc"
)

type Clients struct {
	authConn         *grpc.ClientConn
	courseConn       *grpc.ClientConn
	learningConn     *grpc.ClientConn
	notificationConn *grpc.ClientConn

	Auth          authpb.AuthServiceClient
	Courses       coursepb.CourseServiceClient
	Learning      learningpb.LearningServiceClient
	Notifications notificationpb.NotificationServiceClient
}

func NewClients() (*Clients, error) {
	authConn, err := dial(env("AUTH_SERVICE_ADDR", "localhost:50051"))
	if err != nil {
		return nil, err
	}
	courseConn, err := dial(env("COURSE_SERVICE_ADDR", "localhost:50053"))
	if err != nil {
		_ = authConn.Close()
		return nil, err
	}
	learningConn, err := dial(env("LEARNING_SERVICE_ADDR", "localhost:50054"))
	if err != nil {
		_ = authConn.Close()
		_ = courseConn.Close()
		return nil, err
	}
	notificationConn, err := dial(env("NOTIFICATION_SERVICE_ADDR", "localhost:50055"))
	if err != nil {
		_ = authConn.Close()
		_ = courseConn.Close()
		_ = learningConn.Close()
		return nil, err
	}
	return &Clients{
		authConn: authConn, courseConn: courseConn, learningConn: learningConn, notificationConn: notificationConn,
		Auth: authpb.NewAuthServiceClient(authConn), Courses: coursepb.NewCourseServiceClient(courseConn),
		Learning: learningpb.NewLearningServiceClient(learningConn), Notifications: notificationpb.NewNotificationServiceClient(notificationConn),
	}, nil
}

func (c *Clients) Close() {
	if c == nil {
		return
	}
	_ = c.authConn.Close()
	_ = c.courseConn.Close()
	_ = c.learningConn.Close()
	_ = c.notificationConn.Close()
}

func dial(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
