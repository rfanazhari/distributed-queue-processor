package main

import (
	"context"
	httpHandler "github.com/rfanazhari/distributed-queue-processor/internal/delivery/http"
	rabbitmqRepo "github.com/rfanazhari/distributed-queue-processor/internal/infrastructure/rabbitmq"
	redisRepo "github.com/rfanazhari/distributed-queue-processor/internal/infrastructure/redis"
	usecase2 "github.com/rfanazhari/distributed-queue-processor/internal/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// Default configuration values
	defaultRedisAddr       = "localhost:6379"
	defaultRabbitMQURL     = "amqp://guest:guest@localhost:5672/"
	defaultHTTPAddr        = ":8080"
	defaultShutdownTimeout = 10 * time.Second
)

func main() {
	// Create a context that will be canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChan
		log.Printf("Received signal: %v, initiating shutdown", sig)
		cancel()
	}()

	// Initialize Redis client
	redisAddr := getEnv("REDIS_ADDR", defaultRedisAddr)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer redisClient.Close()

	// Ping Redis to check connection
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Printf("Connected to Redis at %s", redisAddr)

	// Initialize Redis lock repository
	lockRepo := redisRepo.NewLockRepository(redisClient)

	// Initialize RabbitMQ repository
	rabbitMQURL := getEnv("RABBITMQ_URL", defaultRabbitMQURL)
	queueRepo, err := rabbitmqRepo.NewQueueRepository(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer queueRepo.Close()
	log.Printf("Connected to RabbitMQ at %s", rabbitMQURL)

	// Initialize use cases
	reportUseCase := usecase2.NewReportUseCase(lockRepo, queueRepo)
	reportWorkerUseCase := usecase2.NewReportWorkerUseCase(queueRepo, reportUseCase)

	// Start the worker
	if err := reportWorkerUseCase.StartWorker(ctx); err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}
	log.Println("Report worker started")

	// Initialize HTTP handler
	handler := httpHandler.NewHandler(reportUseCase)
	router := handler.SetupRoutes()

	// Create HTTP server
	httpAddr := getEnv("HTTP_ADDR", defaultHTTPAddr)
	server := &http.Server{
		Addr:    httpAddr,
		Handler: router,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("HTTP server listening on %s", httpAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for context cancellation (from signal handler)
	<-ctx.Done()
	log.Println("Shutting down...")

	// Stop the worker
	if err := reportWorkerUseCase.StopWorker(); err != nil {
		log.Printf("Error stopping worker: %v", err)
	}

	// Gracefully shut down the HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Shutdown complete")
}

// getEnv gets an environment variable or returns the default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
