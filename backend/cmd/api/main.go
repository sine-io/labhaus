package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labhaus/backend/internal/application/command"
	"github.com/labhaus/backend/internal/application/query"
	"github.com/labhaus/backend/internal/infrastructure/auth"
	"github.com/labhaus/backend/internal/infrastructure/config"
	httpInfra "github.com/labhaus/backend/internal/infrastructure/http"
	"github.com/labhaus/backend/internal/infrastructure/http/handlers"
	"github.com/labhaus/backend/internal/infrastructure/logger"
	"github.com/labhaus/backend/internal/infrastructure/persistence"
	"github.com/labhaus/backend/internal/infrastructure/queue"
	"github.com/redis/go-redis/v9"
)

const version = "0.1.0"

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(cfg.Log.Level, cfg.Log.Format)
	logger.SetGlobalLogger(log)

	log.Info("Starting Labhaus API", "version", version, "env", cfg.Server.Environment)

	// Initialize database
	db, err := persistence.NewDB(persistence.DBConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	log.Info("Connected to database", "host", cfg.Database.Host, "db", cfg.Database.DBName)

	// Run migrations
	if err := persistence.AutoMigrate(db); err != nil {
		log.Fatal("Failed to run migrations", err)
	}

	log.Info("Database migrations completed")

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test Redis connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis", err)
	}

	log.Info("Connected to Redis", "host", cfg.Redis.Host, "db", cfg.Redis.DB)

	// Initialize queue
	taskQueue := queue.NewRedisQueue(redisClient, "labhaus")

	// Start queue worker
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	taskQueue.StartWorker(workerCtx, queue.WorkflowTaskHandler())
	log.Info("Queue worker started")

	// Initialize password hasher
	passwordHasher := persistence.NewBcryptHasher()

	// Initialize JWT service
	jwtService := auth.NewJWTService(
		cfg.JWT.SecretKey,
		time.Duration(cfg.JWT.TokenDuration)*time.Hour,
	)

	log.Info("JWT service initialized", "token_duration_hours", cfg.JWT.TokenDuration)

	// Initialize repositories
	styleRepo := persistence.NewStyleRepository(db)
	userRepo := persistence.NewUserRepository(db)
	workflowRepo := persistence.NewWorkflowRepository(db)

	// Initialize application services
	styleQueryHandler := query.NewStyleQueryHandler(styleRepo)
	styleCommandHandler := command.NewStyleCommandHandler(styleRepo)

	userQueryHandler := query.NewUserQueryHandler(userRepo, passwordHasher)
	userCommandHandler := command.NewUserCommandHandler(userRepo, passwordHasher)

	workflowQueryHandler := query.NewWorkflowQueryHandler(workflowRepo)
	workflowCommandHandler := command.NewWorkflowCommandHandler(workflowRepo)

	// Initialize HTTP handlers
	healthHandler := handlers.NewHealthHandler(version)
	styleHandler := handlers.NewStyleHandler(styleQueryHandler, styleCommandHandler)
	userHandler := handlers.NewUserHandler(userQueryHandler, userCommandHandler, jwtService)
	workflowHandler := handlers.NewWorkflowHandler(workflowQueryHandler, workflowCommandHandler)

	// Setup router
	router := httpInfra.NewRouter(healthHandler, styleHandler, userHandler, workflowHandler, jwtService, log)
	router.Setup()

	// Create HTTP server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router.Engine(),
	}

	// Start server in a goroutine
	go func() {
		log.Info("HTTP server listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	// Stop queue worker
	workerCancel()
	if err := taskQueue.Close(); err != nil {
		log.Error("Error closing queue", err)
	}

	// Close Redis connection
	if err := redisClient.Close(); err != nil {
		log.Error("Error closing Redis connection", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", err)
	}

	log.Info("Server exited")
}
