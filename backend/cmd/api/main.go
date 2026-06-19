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
	"github.com/labhaus/backend/internal/infrastructure/config"
	httpInfra "github.com/labhaus/backend/internal/infrastructure/http"
	"github.com/labhaus/backend/internal/infrastructure/http/handlers"
	"github.com/labhaus/backend/internal/infrastructure/logger"
	"github.com/labhaus/backend/internal/infrastructure/persistence"
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

	// Initialize repositories
	styleRepo := persistence.NewStyleRepository(db)
	// userRepo := persistence.NewUserRepository(db)
	// workflowRepo := persistence.NewWorkflowRepository(db)

	// Initialize application services
	styleQueryHandler := query.NewStyleQueryHandler(styleRepo)
	styleCommandHandler := command.NewStyleCommandHandler(styleRepo)

	// Initialize HTTP handlers
	healthHandler := handlers.NewHealthHandler(version)
	styleHandler := handlers.NewStyleHandler(styleQueryHandler, styleCommandHandler)

	// Setup router
	router := httpInfra.NewRouter(healthHandler, styleHandler, log)
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

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", err)
	}

	log.Info("Server exited")
}
