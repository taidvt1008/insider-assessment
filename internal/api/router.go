package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"insider-message-sender/internal/config"
	"insider-message-sender/internal/repository"
	"insider-message-sender/internal/scheduler"

	_ "insider-message-sender/internal/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	httpServer *http.Server
}

// @title Insider Message Sender API
// @version 1.0
// @description Golang-based automatic message sending service
// @host localhost:8080
// @BasePath /
func NewServer(cfg *config.Config, s *scheduler.Scheduler, repo *repository.MessageRepository) *Server {
	r := gin.Default()

	// Health check endpoint (no versioning needed)
	r.GET("/health", HealthCheck(s, repo))

	v1 := r.Group("/api/v1")
	v1.POST("/scheduler/start", StartScheduler(s))
	v1.POST("/scheduler/stop", StopScheduler(s))
	v1.GET("/messages/sent", GetSentMessages(repo))
	v1.GET("/messages/failed", GetFailedMessages(repo))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	httpServer := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	log.Printf("HTTP server listening on port %s", cfg.ServerPort)
	return &Server{httpServer: httpServer}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Println("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}
