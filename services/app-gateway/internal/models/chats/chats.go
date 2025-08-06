package chats

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Chat представляет собой структуру чата в системе.
type Chat struct {
	ID         string `json:"id" db:"id"`
	UserFromID string `json:"user_from_id" db:"user_from_id"`
	UserToID   string `json:"user_to_id" db:"user_to_id"`
	CreatedAt  string `json:"created_at" db:"created_at"`
}

// ChatInterface определяет методы, которые должен реализовать модель чата.
type ChatInterface interface {
	Create(ctx context.Context, pool *pgxpool.Pool) error
	Update(ctx context.Context, pool *pgxpool.Pool) error
	ReadByID(ctx context.Context, pool *pgxpool.Pool) (*Chat, error)
	ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Chat, error)
	Delete(ctx context.Context, pool *pgxpool.Pool) error
}

// Create создает новый чат в базе данных.
func (c *Chat) Create(ctx context.Context, pool *pgxpool.Pool) error {
	query := `INSERT INTO chats (user_from_id, user_to_id, created_at)
		VALUES ($1, $2, $3) RETURNING id;`

	err := pool.QueryRow(ctx, query, c.UserFromID, c.UserToID, c.CreatedAt).Scan(&c.ID)
	if err != nil {
		return err
	}

	return nil
}

// Update обновляет существующий чат в базе данных.
func (c *Chat) Update(ctx context.Context, pool *pgxpool.Pool) error {
	query := `UPDATE chats SET user_from_id = $1, user_to_id = $2, created_at = $3 WHERE id = $4;`

	_, err := pool.Exec(ctx, query, c.UserFromID, c.UserToID, c.CreatedAt, c.ID)
	if err != nil {
		return err
	}

	return nil
}

// ReadByID предоставляет чат по ID из базы данных.
func (c *Chat) ReadByID(ctx context.Context, pool *pgxpool.Pool) (*Chat, error) {
	query := `SELECT id, user_from_id, user_to_id, created_at FROM chats WHERE id = $1;`

	err := pool.QueryRow(ctx, query, c.ID).Scan(&c.ID, &c.UserFromID, &c.UserToID, &c.CreatedAt)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// ReadAll предоставляет все чаты из базы данных с пагинацией.
func (c *Chat) ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Chat, error) {
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

// Delete удаляет чат из базы данных.
func (c *Chat) Delete(ctx context.Context, pool *pgxpool.Pool) error {
	query := `DELETE FROM chats WHERE id = $1;`

	_, err := pool.Exec(ctx, query, c.ID)
	if err != nil {
		return err
	}

	return nil
}
