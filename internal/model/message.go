package model

import "time"

type Message struct {
	ID          int64     `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	Content     string    `json:"content"`
	Status      string    `json:"status"`
	SentAt      time.Time `json:"sent_at"`
}
