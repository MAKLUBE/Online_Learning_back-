package app

import (
	"fmt"
	"log"
	"net"
	"os"

	"online-learning-platform/internal/adapters/db"
	"online-learning-platform/internal/adapters/metrics"
	"online-learning-platform/internal/adapters/natsbus"
	"online-learning-platform/internal/services/notifications/mailer"
	"online-learning-platform/internal/services/notifications/repository"
	"online-learning-platform/internal/services/notifications/repository/memory"
	postgresrepo "online-learning-platform/internal/services/notifications/repository/postgre"
	"online-learning-platform/internal/services/notifications/subscriber"
	notificationgrpc "online-learning-platform/internal/services/notifications/transport/grpc"
	"online-learning-platform/internal/services/notifications/usecase"
	notificationpb "online-learning-platform/pkg/gen/notifications"

	"google.golang.org/grpc"
)

func Run() error {
	port := getEnv("GRPC_PORT", getEnv("PORT", "50055"))
	metrics.Start(getEnv("METRICS_PORT", "9104"), "notification-service")
	var repo repository.Repository = memory.New()
	if database, err := db.ConnectAndMigrate(getEnv("NOTIFICATION_DATABASE_URL", "postgres://postgres:postgres@localhost:5432/notification_db?sslmode=disable"), "migrations/notification-service/000001_create_notifications.up.sql"); err != nil {
		log.Printf("warning: PostgreSQL unavailable, notifications use memory repository: %v", err)
	} else {
		log.Println("PostgreSQL available, notifications use notification_db")
		repo = postgresrepo.New(database)
	}
	uc := usecase.New(repo, mailer.NewSMTPFromEnv())
	eventHandler := subscriber.New(uc)
	if natsSubscriber, err := natsbus.NewSubscriberFromEnv(); err != nil {
		log.Printf("warning: NATS unavailable, notification event subscriber disabled: %v", err)
	} else {
		log.Println("NATS available, notification service subscribes to events")
		defer natsSubscriber.Close()
		for _, subject := range []string{"user.registered", "user.password_reset_requested", "course.enrolled", "course.completed", "assignment.graded"} {
			if err := natsSubscriber.Subscribe(subject, eventHandler.Handle); err != nil {
				log.Printf("warning: subscribe %s failed: %v", subject, err)
			}
		}
	}
	handler := notificationgrpc.NewHandler(uc)
	_ = handler

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("listen on port %s: %w", port, err)
	}
	server := grpc.NewServer()
	notificationpb.RegisterNotificationServiceServer(server, handler)
	return server.Serve(listener)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
