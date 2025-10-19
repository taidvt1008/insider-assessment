package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"insider-message-sender/internal/cache"
	"insider-message-sender/internal/config"
	"insider-message-sender/internal/model"
	"insider-message-sender/internal/repository"
)

type Scheduler struct {
	cfg       *config.Config
	repo      *repository.MessageRepository
	cache     *cache.RedisClient
	client    *http.Client
	isRunning bool
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.Mutex
}

func NewScheduler(cfg *config.Config, repo *repository.MessageRepository, cache *cache.RedisClient) *Scheduler {
	return &Scheduler{
		cfg:   cfg,
		repo:  repo,
		cache: cache,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSHandshakeTimeout:   5 * time.Second,
				ResponseHeaderTimeout: 5 * time.Second,
				MaxIdleConns:          10,
				MaxIdleConnsPerHost:   2,
				IdleConnTimeout:       cfg.SendInterval + 30*time.Second,
			},
		},
	}
}

func (s *Scheduler) Start() error {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return nil
	}

	// Create new context for this start cycle
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.isRunning = true
	s.mu.Unlock()

	log.Println("Scheduler started...")
	go s.run()

	return nil
}

func (s *Scheduler) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	s.cancel() // Cancel context to stop all ongoing operations
	s.isRunning = false

	log.Println("Scheduler stopped!!!")
	return nil
}

func (s *Scheduler) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isRunning
}

func (s *Scheduler) run() {
	s.process()

	ticker := time.NewTicker(s.cfg.SendInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.process()
		case <-s.ctx.Done():
			log.Println("Scheduler context cancelled, stopping...")
			return
		}
	}
}

func (s *Scheduler) process() {
	msgs, err := s.repo.FetchUnsent(2)
	if err != nil {
		log.Printf("DB fetch error: %v", err)
		return
	}

	log.Printf("Fetched %d unsent messages", len(msgs))

	var wg sync.WaitGroup
	for _, m := range msgs {
		wg.Add(1)
		go func(msg model.Message) {
			defer wg.Done()
			s.sendMessage(s.ctx, msg)
		}(m)
	}
	wg.Wait()
}

const (
	maxMessageLength = 160
	redisKeyPrefix   = "insider:msg:sent"
	maxRetries       = 3
	baseDelay        = 1 * time.Second
)

func (s *Scheduler) sendMessage(ctx context.Context, m model.Message) {
	if len(m.Content) > maxMessageLength {
		log.Printf("Message %d content too long (%d chars), skipping", m.ID, len(m.Content))
		return
	}

	body, err := json.Marshal(map[string]string{
		"to":      m.PhoneNumber,
		"content": m.Content,
	})
	if err != nil {
		log.Printf("Failed to marshal message %d: %v", m.ID, err)
		return
	}

	// Retry mechanism for rate limiting and temporary failures
	for attempt := 0; attempt < maxRetries; attempt++ {
		success := s.sendMessageWithRetry(ctx, m, body, attempt)
		if success {
			return
		}

		// Don't retry on last attempt
		if attempt == maxRetries-1 {
			log.Printf("Message %d failed after %d attempts, marking as failed", m.ID, maxRetries)
			if err := s.repo.MarkAsFailed(m.ID); err != nil {
				log.Printf("Failed to mark msg %d as failed in DB: %v", m.ID, err)
			}
			return
		}

		// Calculate delay with exponential backoff
		delay := baseDelay * time.Duration(1<<attempt) // 1s, 2s, 4s
		log.Printf("Message %d attempt %d failed, retrying in %v", m.ID, attempt+1, delay)

		select {
		case <-ctx.Done():
			log.Printf("Message %d retry cancelled due to context cancellation", m.ID)
			return
		case <-time.After(delay):
			// Continue to next attempt
		}
	}
}

func (s *Scheduler) sendMessageWithRetry(ctx context.Context, m model.Message, body []byte, attempt int) bool {
	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "POST", s.cfg.WebhookURL, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to create request for msg %d (attempt %d): %v", m.ID, attempt+1, err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("Failed to send msg %d (attempt %d): %v", m.ID, attempt+1, err)
		return false
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusAccepted {
		log.Printf("Message %d sent successfully (attempt %d)", m.ID, attempt+1)

		var respData struct {
			MessageID string `json:"messageId"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			log.Printf("Failed to parse webhook response for msg %d (attempt %d): %v", m.ID, attempt+1, err)
			return false
		}

		// Use the same timestamp for both DB and cache
		sentAt := time.Now()

		// Mark DB as sent
		if err := s.repo.MarkAsSent(m.ID); err != nil {
			log.Printf("Failed to mark msg %d as sent in DB: %v", m.ID, err)
		}

		// Cache messageId + sending time
		if respData.MessageID != "" {
			cacheKey := fmt.Sprintf("%s:%s", redisKeyPrefix, respData.MessageID)
			cacheVal := sentAt.Format(time.RFC3339)

			if err := s.cache.Set(ctx, cacheKey, cacheVal, 0); err != nil {
				log.Printf("Failed to cache messageId %s: %v", respData.MessageID, err)
			} else {
				log.Printf("Cached messageId=%s sent_at=%s", respData.MessageID, cacheVal)
			}
		}
		return true
	} else {
		log.Printf("Failed to send msg %d (attempt %d): %s", m.ID, attempt+1, resp.Status)
		// Read response body to avoid connection leak
		_, _ = io.Copy(io.Discard, resp.Body)
		return false // Will retry
	}
}
