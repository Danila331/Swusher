package chats

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Mesage структура представляет собой сообщение в чате.
type Message struct {
	ID        string `json:"id" db:"id"`
	ChatID    string `json:"chat_id" db:"chat_id"`
	UserID    string `json:"user_id" db:"user_id"`
	Text      string `json:"text" db:"text"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

// MessageInterface определяет методы, которые должен реализовать модель сообщения.
type MessageInterface interface {
	Create(ctx context.Context, pool *pgxpool.Pool) error
	Update(ctx context.Context, pool *pgxpool.Pool) error
	ReadByID(ctx context.Context, pool *pgxpool.Pool) (*Message, error)
	ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Message, error)
	Delete(ctx context.Context, pool *pgxpool.Pool) error
}

// Create создает новое сообщение в базе данных.
func (m *Message) Create(ctx context.Context, pool *pgxpool.Pool) error {
	query := `INSERT INTO messages (chat_id, user_id, text, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id;`

	err := pool.QueryRow(ctx, query, m.ChatID, m.UserID, m.Text, m.CreatedAt).Scan(&m.ID)
	if err != nil {
		return err
	}

	return nil
}

// Update обновляет существующее сообщение в базе данных.
func (m *Message) Update(ctx context.Context, pool *pgxpool.Pool) error {
	query := `UPDATE messages SET chat_id = $1, user_id = $2, text = $3, created_at = $4 WHERE id = $5;`

	_, err := pool.Exec(ctx, query, m.ChatID, m.UserID, m.Text, m.CreatedAt, m.ID)
	if err != nil {
		return err
	}

	return nil
}

// ReadByID предоставляет сообщение по ID из базы данных.
func (m *Message) ReadByID(ctx context.Context, pool *pgxpool.Pool) (*Message, error) {
	query := `SELECT chat_id, user_id, text, created_at FROM messages WHERE id = $1;`

	err := pool.QueryRow(ctx, query, m.ID).Scan(&m.ChatID, &m.UserID, &m.Text, &m.CreatedAt)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// ReadAll предоставляет все сообщения из базы данных с пагинацией.
func (m *Message) ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Message, error) {
	query := `SELECT chat_id, user_id, text, created_at FROM messages LIMIT $1 OFFSET $2;`

	rows, err := pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ChatID, &message.UserID, &message.Text, &message.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// Delete удаляет сообщение из базы данных.
func (m *Message) Delete(ctx context.Context, pool *pgxpool.Pool) error {
	query := `DELETE FROM messages WHERE id = $1;`

	_, err := pool.Exec(ctx, query, m.ID)
	if err != nil {
		return err
	}

	return nil
}
