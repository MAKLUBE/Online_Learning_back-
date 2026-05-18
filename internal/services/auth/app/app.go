package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"online-learning-platform/internal/adapters/events"
	"online-learning-platform/internal/adapters/metrics"
	"online-learning-platform/internal/adapters/natsbus"
	platformredis "online-learning-platform/internal/adapters/redis"
	"online-learning-platform/internal/services/auth/repository"
	"online-learning-platform/internal/services/auth/repository/memory"
	postgresrepo "online-learning-platform/internal/services/auth/repository/postgre"
	redisrepo "online-learning-platform/internal/services/auth/repository/redis"
	authgrpc "online-learning-platform/internal/services/auth/transport/grpc"
	"online-learning-platform/internal/services/auth/usecase"
	authpb "online-learning-platform/pkg/gen/auth"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func Run() error {
	port := getEnv("GRPC_PORT", getEnv("PORT", "50051"))
	metrics.Start(getEnv("METRICS_PORT", "9101"), "user-auth-service")

	var sessionRepo repository.AuthRepository = memory.NewAuthRepository()
	var userRepo repository.UserRepository = memory.NewAuthRepository()

	if db, err := connectAuthDB(); err != nil {
		log.Printf("warning: PostgreSQL unavailable, auth users use memory repository: %v", err)
	} else {
		log.Println("PostgreSQL available, auth users use user_auth_db")
		pgRepo := postgresrepo.New(db)
		userRepo = pgRepo.(repository.UserRepository)
	}

	redisClient := platformredis.NewFromEnv()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx); err != nil {
		log.Printf("warning: Redis unavailable, auth sessions use memory repository: %v", err)
	} else {
		log.Println("Redis available, auth refresh sessions use Redis")
		sessionRepo = redisrepo.New(redisClient)
	}

	var publisher events.Publisher = events.LogPublisher{}
	if natsPublisher, err := natsbus.NewPublisherFromEnv(); err != nil {
		log.Printf("warning: NATS unavailable, auth events use log publisher: %v", err)
	} else {
		log.Println("NATS available, auth events publish to NATS")
		defer natsPublisher.Close()
		publisher = natsPublisher
	}

	uc := usecase.NewAuthUseCase(sessionRepo)
	handler := authgrpc.NewHandlerWithDependencies(uc, userRepo, publisher)
	return serve(port, handler)
}

func serve(port string, handler authpb.AuthServiceServer) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("listen on port %s: %w", port, err)
	}

	server := grpc.NewServer()
	authpb.RegisterAuthServiceServer(server, handler)

	return server.Serve(listener)
}

func connectAuthDB() (*sql.DB, error) {
	dsn := getEnv("AUTH_DATABASE_URL", "postgres://postgres:postgres@localhost:5432/user_auth_db?sslmode=disable")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := runAuthMigration(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func runAuthMigration(ctx context.Context, db *sql.DB) error {
	path := filepath.Join("migrations", "user-service", "000001_create_users.up.sql")
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, string(content))
	return err
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
