package model

import "time"

type SentMessageResponseData struct {
	ID          int64     `json:"id" example:"1"`
	PhoneNumber string    `json:"phone_number" example:"+84901234567"`
	Content     string    `json:"content" example:"Hello from Insider!"`
	Status      string    `json:"status" example:"sent"`
	SentAt      time.Time `json:"sent_at" example:"2025-10-19T07:41:45Z"`
}

type Pagination struct {
	Limit   int  `json:"limit" example:"10"`
	Offset  int  `json:"offset" example:"0"`
	Count   int  `json:"count" example:"2"`
	Total   int  `json:"total" example:"5"`
	HasMore bool `json:"has_more" example:"true"`
}

type SentMessagesResponse struct {
	Data       []SentMessageResponseData `json:"data"`
	Pagination Pagination                `json:"pagination"`
}

type SchedulerActionResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Scheduler started successfully"`
	Time    string `json:"time" example:"2025-10-19T08:10:00Z"`
}

type ErrorResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Internal server error"`
	Time    string `json:"time" example:"2025-10-19T09:00:00Z"`
}

type HealthResponse struct {
	Status    string            `json:"status" example:"healthy"`
	Timestamp string            `json:"timestamp" example:"2025-10-19T09:00:00Z"`
	Services  map[string]string `json:"services" example:"database:healthy,scheduler:running,redis:healthy"`
}
