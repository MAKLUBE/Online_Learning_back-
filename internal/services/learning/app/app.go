package app

import (
	"fmt"
	"log"
	"net"
	"os"

	"online-learning-platform/internal/adapters/db"
	"online-learning-platform/internal/adapters/events"
	"online-learning-platform/internal/adapters/metrics"
	"online-learning-platform/internal/adapters/natsbus"
	"online-learning-platform/internal/services/learning/repository"
	"online-learning-platform/internal/services/learning/repository/memory"
	postgresrepo "online-learning-platform/internal/services/learning/repository/postgre"
	"online-learning-platform/internal/services/learning/transport/client"
	learninggrpc "online-learning-platform/internal/services/learning/transport/grpc"
	"online-learning-platform/internal/services/learning/usecase"
	learningpb "online-learning-platform/pkg/gen/learning"

	"google.golang.org/grpc"
)

func Run() error {
	port := getEnv("GRPC_PORT", getEnv("PORT", "50054"))
	metrics.Start(getEnv("METRICS_PORT", "9103"), "learning-service")
	var repo repository.Repository = memory.New()
	if database, err := db.ConnectAndMigrate(getEnv("LEARNING_DATABASE_URL", "postgres://postgres:postgres@localhost:5432/learning_db?sslmode=disable"), "migrations/learning-service/000001_create_learning.up.sql"); err != nil {
		log.Printf("warning: PostgreSQL unavailable, learning uses memory repository: %v", err)
	} else {
		log.Println("PostgreSQL available, learning uses learning_db")
		repo = postgresrepo.New(database)
	}
	var userChecker usecase.UserChecker = usecase.AlwaysExistsChecker{}
	if authChecker, err := client.NewAuthUserChecker(getEnv("AUTH_SERVICE_ADDR", "localhost:50051")); err != nil {
		log.Printf("warning: auth checker unavailable, learning user checks are permissive: %v", err)
	} else {
		defer authChecker.Close()
		userChecker = authChecker
	}
	var courseChecker usecase.CourseChecker = usecase.AlwaysExistsChecker{}
	if checker, err := client.NewCourseChecker(getEnv("COURSE_SERVICE_ADDR", "localhost:50053")); err != nil {
		log.Printf("warning: course checker unavailable, learning course checks are permissive: %v", err)
	} else {
		defer checker.Close()
		courseChecker = checker
	}
	var publisher events.Publisher = events.LogPublisher{}
	if natsPublisher, err := natsbus.NewPublisherFromEnv(); err != nil {
		log.Printf("warning: NATS unavailable, learning events use log publisher: %v", err)
	} else {
		log.Println("NATS available, learning events publish to NATS")
		defer natsPublisher.Close()
		publisher = natsPublisher
	}
	uc := usecase.New(repo, userChecker, courseChecker, publisher)
	handler := learninggrpc.NewHandler(uc)
	_ = handler

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("listen on port %s: %w", port, err)
	}
	server := grpc.NewServer()
	learningpb.RegisterLearningServiceServer(server, handler)
	return server.Serve(listener)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
