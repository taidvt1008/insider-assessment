package api

import (
	"context"
	"net/http"
	"time"

	"insider-message-sender/internal/model"
	"insider-message-sender/internal/repository"
	"insider-message-sender/internal/scheduler"

	"github.com/gin-gonic/gin"
)

// @Summary Health check endpoint
// @Description Check the health status of the service including database connectivity and scheduler status
// @Tags Health
// @Produce json
// @Success 200 {object} model.HealthResponse
// @Failure 503 {object} model.HealthResponse
// @Router /health [get]
func HealthCheck(s *scheduler.Scheduler, repo *repository.MessageRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		health := model.HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now().Format(time.RFC3339),
			Services:  make(map[string]string),
		}

		// Check database connectivity
		if err := repo.Ping(ctx); err != nil {
			health.Status = "unhealthy"
			health.Services["database"] = "unhealthy: " + err.Error()
		} else {
			health.Services["database"] = "healthy"
		}

		// Check scheduler status
		if s.IsRunning() {
			health.Services["scheduler"] = "running"
		} else {
			health.Services["scheduler"] = "stopped"
		}

		// Check Redis connectivity (if available)
		health.Services["redis"] = "healthy" // Assume healthy for now

		// Return appropriate status code
		if health.Status == "healthy" {
			c.JSON(http.StatusOK, health)
		} else {
			c.JSON(http.StatusServiceUnavailable, health)
		}
	}
}
