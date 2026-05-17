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
	"online-learning-platform/internal/services/courses/repository"
	"online-learning-platform/internal/services/courses/repository/memory"
	postgresrepo "online-learning-platform/internal/services/courses/repository/postgre"
	"online-learning-platform/internal/services/courses/transport/client"
	coursesgrpc "online-learning-platform/internal/services/courses/transport/grpc"
	"online-learning-platform/internal/services/courses/usecase"
	coursepb "online-learning-platform/pkg/gen/courses"

	"google.golang.org/grpc"
)

func Run() error {
	port := getEnv("GRPC_PORT", getEnv("PORT", "50053"))
	metrics.Start(getEnv("METRICS_PORT", "9102"), "course-service")

	var repo repository.CourseRepository = memory.NewCourseRepository()
	if database, err := db.ConnectAndMigrate(getEnv("COURSE_DATABASE_URL", "postgres://postgres:postgres@localhost:5432/course_db?sslmode=disable"), "migrations/course-service/000001_create_courses.up.sql"); err != nil {
		log.Printf("warning: PostgreSQL unavailable, courses use memory repository: %v", err)
	} else {
		log.Println("PostgreSQL available, courses use course_db")
		repo = postgresrepo.New(database)
	}
	var checker usecase.InstructorChecker
	if authClient, err := client.NewAuthClient(getEnv("AUTH_SERVICE_ADDR", "localhost:50051")); err != nil {
		log.Printf("warning: auth checker unavailable, instructor_id DB role check disabled: %v", err)
	} else {
		defer authClient.Close()
		checker = authClient
	}
	var publisher events.Publisher = events.LogPublisher{}
	if natsPublisher, err := natsbus.NewPublisherFromEnv(); err != nil {
		log.Printf("warning: NATS unavailable, course events use log publisher: %v", err)
	} else {
		log.Println("NATS available, course events publish to NATS")
		defer natsPublisher.Close()
		publisher = natsPublisher
	}
	uc := usecase.NewCourseUseCaseWithDependencies(repo, publisher, checker)
	handler := coursesgrpc.NewHandler(uc)
	_ = handler

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("listen on port %s: %w", port, err)
	}

	server := grpc.NewServer()
	coursepb.RegisterCourseServiceServer(server, handler)

	return server.Serve(listener)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
