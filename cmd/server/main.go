package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"insider-message-sender/internal/api"
	"insider-message-sender/internal/cache"
	"insider-message-sender/internal/config"
	"insider-message-sender/internal/repository"
	"insider-message-sender/internal/scheduler"
)

func main() {
	cfg := config.Load()
	log.Printf("DB Host: %s, Redis Host: %s, WebhookURL: %s, SendInterval: %s, ServerPort: %s",
		cfg.DBHost, cfg.RedisHost, cfg.WebhookURL, cfg.SendInterval, cfg.ServerPort)

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	repo := repository.NewMessageRepository(connStr)
	// Ensure database connection is closed on exit
	defer func() {
		log.Println("Closing database connection...")
		if err := repo.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	redisClient := cache.NewRedisClient(cfg.RedisHost)
	// Ensure Redis connection is closed on exit
	defer func() {
		log.Println("Closing Redis connection...")
		if err := redisClient.Close(); err != nil {
			log.Printf("Error closing Redis: %v", err)
		}
	}()

	s := scheduler.NewScheduler(cfg, repo, redisClient)
	if err := s.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create HTTP server
	server := api.NewServer(cfg, s, repo)

	// Start server in a goroutine with error handling
	serverErr := make(chan error, 1)
	go func() {
		log.Println("Starting HTTP server...")
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// Wait for either shutdown signal or server error
	select {
	case err := <-serverErr:
		log.Printf("HTTP server failed: %v", err)
		log.Println("Stopping scheduler due to server failure...")
		if stopErr := s.Stop(); stopErr != nil {
			log.Printf("Error stopping scheduler: %v", stopErr)
		}
		log.Fatal("Application terminated due to server failure")
	case <-sigChan:
		log.Println("Received shutdown signal, starting graceful shutdown...")

		// Stop scheduler first
		if err := s.Stop(); err != nil {
			log.Printf("Error stopping scheduler: %v", err)
		}

		// Gracefully shutdown HTTP server
		if err := server.Shutdown(30 * time.Second); err != nil {
			log.Printf("Error shutting down HTTP server: %v", err)
		}

		log.Println("Application shutdown complete.")
	}
}
