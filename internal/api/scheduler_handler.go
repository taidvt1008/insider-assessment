package api

import (
	"net/http"
	"time"

	"insider-message-sender/internal/model"
	"insider-message-sender/internal/scheduler"

	"github.com/gin-gonic/gin"
)

// @Summary Start automatic message sending
// @Description Starts the background scheduler that periodically sends pending messages every configured interval.
// @Tags Scheduler
// @Produce json
// @Success 200 {object} model.SchedulerActionResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/scheduler/start [post]
func StartScheduler(s *scheduler.Scheduler) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := s.Start()
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Status:  "error",
				Message: "Internal server error",
				Time:    time.Now().Format(time.RFC3339),
			})
			return
		}
		c.JSON(http.StatusOK, model.SchedulerActionResponse{
			Status:  "success",
			Message: "Scheduler started successfully",
			Time:    time.Now().Format(time.RFC3339),
		})
	}
}

// @Summary Stop automatic message sending
// @Description Stops the background scheduler. No further messages will be sent until restarted.
// @Tags Scheduler
// @Produce json
// @Success 200 {object} model.SchedulerActionResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/scheduler/stop [post]
func StopScheduler(s *scheduler.Scheduler) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := s.Stop()
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Status:  "error",
				Message: "Internal server error",
				Time:    time.Now().Format(time.RFC3339),
			})
			return
		}
		c.JSON(http.StatusOK, model.SchedulerActionResponse{
			Status:  "success",
			Message: "Scheduler stopped successfully",
			Time:    time.Now().Format(time.RFC3339),
		})
	}
}
