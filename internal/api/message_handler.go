package api

import (
	"net/http"
	"strconv"

	"insider-message-sender/internal/model"
	"insider-message-sender/internal/repository"

	"github.com/gin-gonic/gin"
)

// @Summary Get list of sent messages (with pagination)
// @Tags Messages
// @Produce json
// @Param limit query int false "Number of messages to return" default(10)
// @Param offset query int false "Number of messages to skip" default(0)
// @Success 200 {object} model.SentMessagesResponse
// @Router /api/v1/messages/sent [get]
func GetSentMessages(repo *repository.MessageRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil || limit <= 0 {
			limit = 10
		}

		offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
		if err != nil || offset < 0 {
			offset = 0
		}

		msgs, err := repo.FetchSent(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		total, err := repo.CountSent()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := model.SentMessagesResponse{
			Data: make([]model.SentMessageResponseData, len(msgs)),
			Pagination: model.Pagination{
				Limit:   limit,
				Offset:  offset,
				Count:   len(msgs),
				Total:   total,
				HasMore: offset+limit < total,
			},
		}

		for i, m := range msgs {
			resp.Data[i] = model.SentMessageResponseData{
				ID:          m.ID,
				PhoneNumber: m.PhoneNumber,
				Content:     m.Content,
				Status:      m.Status,
				SentAt:      m.SentAt,
			}
		}

		c.JSON(http.StatusOK, resp)
	}
}

// @Summary Get list of failed messages (with pagination)
// @Tags Messages
// @Produce json
// @Param limit query int false "Number of messages to return" default(10)
// @Param offset query int false "Number of messages to skip" default(0)
// @Success 200 {object} model.SentMessagesResponse
// @Router /api/v1/messages/failed [get]
func GetFailedMessages(repo *repository.MessageRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil || limit <= 0 {
			limit = 10
		}

		offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
		if err != nil || offset < 0 {
			offset = 0
		}

		msgs, err := repo.FetchFailed(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		total, err := repo.CountFailed()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := model.SentMessagesResponse{
			Data: make([]model.SentMessageResponseData, len(msgs)),
			Pagination: model.Pagination{
				Limit:   limit,
				Offset:  offset,
				Count:   len(msgs),
				Total:   total,
				HasMore: offset+limit < total,
			},
		}

		for i, m := range msgs {
			resp.Data[i] = model.SentMessageResponseData{
				ID:          m.ID,
				PhoneNumber: m.PhoneNumber,
				Content:     m.Content,
				Status:      m.Status,
				SentAt:      m.SentAt,
			}
		}

		c.JSON(http.StatusOK, resp)
	}
}
