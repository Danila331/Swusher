package chats

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Chat struct represents a chat in the system.
type Chat struct {
	ID         string `json:"id" db:"id"`
	UserFromID string `json:"user_from_id" db:"user_from_id"`
	UserToID   string `json:"user_to_id" db:"user_to_id"`
	CreatedAt  string `json:"created_at" db:"created_at"`
}

// ChatInterface defines the methods that a chat model should implement.
type ChatInterface interface {
	Create(ctx context.Context, pool *pgxpool.Pool) error
	Update(ctx context.Context, pool *pgxpool.Pool) error
	ReadByID(ctx context.Context, pool *pgxpool.Pool) (*Chat, error)
	ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Chat, error)
	Delete(ctx context.Context, pool *pgxpool.Pool) error
}

// Create creates a new chat in the database.
func (c *Chat) Create(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to create a new chat in the database
	// This is just a placeholder implementation

	query := `INSERT INTO chats (user_from_id, user_to_id, created_at)
		VALUES ($1, $2, $3) RETURNING id;`

	err := pool.QueryRow(ctx, query, c.UserFromID, c.UserToID, c.CreatedAt).Scan(&c.ID)
	if err != nil {
		return err
	}

	return nil
}

// Update updates an existing chat in the database.
func (c *Chat) Update(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to update an existing chat in the database
	// This is just a placeholder implementation

	query := `UPDATE chats SET user_from_id = $1, user_to_id = $2, created_at = $3 WHERE id = $4;`

	_, err := pool.Exec(ctx, query, c.UserFromID, c.UserToID, c.CreatedAt, c.ID)
	if err != nil {
		return err
	}

	return nil
}

// ReadByID retrieves a chat by its ID from the database.
func (c *Chat) ReadByID(ctx context.Context, pool *pgxpool.Pool) (*Chat, error) {
	// Implement the logic to read a chat by its ID from the database
	// This is just a placeholder implementation

	query := `SELECT id, user_from_id, user_to_id, created_at FROM chats WHERE id = $1;`

	err := pool.QueryRow(ctx, query, c.ID).Scan(&c.ID, &c.UserFromID, &c.UserToID, &c.CreatedAt)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// ReadAll retrieves all chats from the database with pagination.
func (c *Chat) ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Chat, error) {
	// Implement the logic to read all chats from the database with pagination
	// This is just a placeholder implementation

	query := `SELECT id, user_from_id, user_to_id, created_at FROM chats LIMIT $1 OFFSET $2;`

	rows, err := pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.ID, &chat.UserFromID, &chat.UserToID, &chat.CreatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

// Delete removes a chat from the database.
func (c *Chat) Delete(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to delete a chat from the database
	// This is just a placeholder implementation

	query := `DELETE FROM chats WHERE id = $1;`

	_, err := pool.Exec(ctx, query, c.ID)
	if err != nil {
		return err
	}

	return nil
}
