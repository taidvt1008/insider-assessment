package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"insider-message-sender/internal/constants"
	"insider-message-sender/internal/model"

	_ "github.com/lib/pq"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(connStr string) *MessageRepository {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}
	setupPool(db)
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}
	return &MessageRepository{db: db}
}

func setupPool(db *sql.DB) {
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(10 * time.Minute)
}

func (r *MessageRepository) FetchUnsent(limit int) ([]model.Message, error) {
	rows, err := r.db.Query(`SELECT id, phone_number, content, status FROM messages WHERE status = $1 LIMIT $2`, constants.MessageStatusPending, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var msgs []model.Message
	for rows.Next() {
		var m model.Message
		if err := rows.Scan(&m.ID, &m.PhoneNumber, &m.Content, &m.Status); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

func (r *MessageRepository) MarkAsSent(id int64) error {
	_, err := r.db.Exec(`UPDATE messages SET status=$1, sent_at=$2 WHERE id=$3`, constants.MessageStatusSent, time.Now(), id)
	return err
}

func (r *MessageRepository) MarkAsFailed(id int64) error {
	_, err := r.db.Exec(`UPDATE messages SET status=$1, sent_at=$2 WHERE id=$3`, constants.MessageStatusFailed, time.Now(), id)
	return err
}

func (r *MessageRepository) FetchSent(limit, offset int) ([]model.Message, error) {
	query := `SELECT id, phone_number, content, status, sent_at 
			  FROM messages
			  WHERE status = $1
			  ORDER BY sent_at DESC
			  LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, constants.MessageStatusSent, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var msgs []model.Message
	for rows.Next() {
		var m model.Message
		if err := rows.Scan(
			&m.ID,
			&m.PhoneNumber,
			&m.Content,
			&m.Status,
			&m.SentAt,
		); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}

	return msgs, nil
}

func (r *MessageRepository) CountSent() (int, error) {
	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM messages WHERE status = $1`, constants.MessageStatusSent).Scan(&total)
	return total, err
}

func (r *MessageRepository) CountFailed() (int, error) {
	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM messages WHERE status = $1`, constants.MessageStatusFailed).Scan(&total)
	return total, err
}

func (r *MessageRepository) CountPending() (int, error) {
	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM messages WHERE status = $1`, constants.MessageStatusPending).Scan(&total)
	return total, err
}

func (r *MessageRepository) FetchFailed(limit, offset int) ([]model.Message, error) {
	query := `SELECT id, phone_number, content, status, sent_at 
			  FROM messages
			  WHERE status = $1
			  ORDER BY sent_at DESC
			  LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, constants.MessageStatusFailed, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var msgs []model.Message
	for rows.Next() {
		var m model.Message
		if err := rows.Scan(
			&m.ID,
			&m.PhoneNumber,
			&m.Content,
			&m.Status,
			&m.SentAt,
		); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}

	return msgs, nil
}

func (r *MessageRepository) Close() error {
	return r.db.Close()
}

func (r *MessageRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
