package chats

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Mesage struct represents a message in the system.
type Message struct {
	ID        string `json:"id" db:"id"`
	ChatID    string `json:"chat_id" db:"chat_id"`
	UserID    string `json:"user_id" db:"user_id"`
	Text      string `json:"text" db:"text"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

// MessageInterface defines the methods that a message model should implement.
type MessageInterface interface {
	Create(ctx context.Context, pool *pgxpool.Pool) error
	Update(ctx context.Context, pool *pgxpool.Pool) error
	ReadByID(ctx context.Context, pool *pgxpool.Pool) (*Message, error)
	ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Message, error)
	Delete(ctx context.Context, pool *pgxpool.Pool) error
}

// Create creates a new message in the database.
func (m *Message) Create(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to create a new message in the database
	// This is just a placeholder implementation

	query := `INSERT INTO messages (chat_id, user_id, text, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id;`

	err := pool.QueryRow(ctx, query, m.ChatID, m.UserID, m.Text, m.CreatedAt).Scan(&m.ID)
	if err != nil {
		return err
	}

	return nil
}

// Update updates an existing message in the database.
func (m *Message) Update(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to update an existing message in the database
	// This is just a placeholder implementation

	query := `UPDATE messages SET chat_id = $1, user_id = $2, text = $3, created_at = $4 WHERE id = $5;`

	_, err := pool.Exec(ctx, query, m.ChatID, m.UserID, m.Text, m.CreatedAt, m.ID)
	if err != nil {
		return err
	}

	return nil
}

// ReadByID retrieves a message by its ID from the database.
func (m *Message) ReadByID(ctx context.Context, pool *pgxpool.Pool) (*Message, error) {
	// Implement the logic to read a message by its ID from the database
	// This is just a placeholder implementation

	query := `SELECT chat_id, user_id, text, created_at FROM messages WHERE id = $1;`

	err := pool.QueryRow(ctx, query, m.ID).Scan(&m.ChatID, &m.UserID, &m.Text, &m.CreatedAt)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// ReadAll retrieves all messages from the database with pagination.
func (m *Message) ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Message, error) {
	// Implement the logic to read all messages from the database with pagination
	// This is just a placeholder implementation

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

// Delete removes a message from the database.
func (m *Message) Delete(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to delete a message from the database
	// This is just a placeholder implementation

	query := `DELETE FROM messages WHERE id = $1;`

	_, err := pool.Exec(ctx, query, m.ID)
	if err != nil {
		return err
	}

	return nil
}
