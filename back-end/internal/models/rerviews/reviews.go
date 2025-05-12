package rerviews

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Review struct represents a review in the system.
type Review struct {
	ID              string  `json:"id" db:"id"`
	AdvertisementID string  `json:"advertisement_id" db:"advertisement_id"`
	UserFromID      string  `json:"user_from_id" db:"user_from_id"`
	Score           float64 `json:"score" db:"score"`
	Text            string  `json:"text" db:"text"`
	CreatedAt       string  `json:"created_at" db:"created_at"`
}

// ReviewInterface defines the methods that a review model should implement.
type ReviewInterface interface {
	Create(ctx context.Context, pool *pgxpool.Pool) error
	Update(ctx context.Context, pool *pgxpool.Pool) error
	ReadByID(ctx context.Context, pool *pgxpool.Pool) (*Review, error)
	ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Review, error)
	Delete(ctx context.Context, pool *pgxpool.Pool) error
}

// Create creates a new review in the database.
func (r *Review) Create(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to create a new review in the database
	// This is just a placeholder implementation

	query := `INSERT INTO reviews (advertisement_id, user_from_id, score, text, created_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	err := pool.QueryRow(ctx, query, r.AdvertisementID, r.UserFromID, r.Score, r.Text, r.CreatedAt).Scan(&r.ID)
	if err != nil {
		return err
	}

	return nil
}

// Update updates an existing review in the database.
func (r *Review) Update(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to update an existing review in the database
	// This is just a placeholder implementation

	query := `UPDATE reviews SET advertisement_id = $1, user_from_id = $2, score = $3, text = $4, created_at = $5 WHERE id = $6;`

	_, err := pool.Exec(ctx, query, r.AdvertisementID, r.UserFromID, r.Score, r.Text, r.CreatedAt, r.ID)
	if err != nil {
		return err
	}

	return nil
}

// ReadByID retrieves a review by its ID from the database.
func (r *Review) ReadByID(ctx context.Context, pool *pgxpool.Pool) (*Review, error) {
	// Implement the logic to read a review by its ID from the database
	// This is just a placeholder implementation

	query := `SELECT advertisement_id, user_from_id, score, text, created_at FROM reviews WHERE id = $1;`

	err := pool.QueryRow(ctx, query, r.ID).Scan(&r.AdvertisementID, &r.UserFromID, &r.Score, &r.Text, &r.CreatedAt)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// ReadAll retrieves all reviews from the database with pagination.
func (r *Review) ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Review, error) {
	// Implement the logic to read all reviews from the database with pagination
	// This is just a placeholder implementation

	query := `SELECT advertisement_id, user_from_id, score, text, created_at FROM reviews LIMIT $1 OFFSET $2;`

	rows, err := pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		err := rows.Scan(&review.AdvertisementID, &review.UserFromID, &review.Score, &review.Text, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

// Delete deletes a review from the database.
func (r *Review) Delete(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to delete a review from the database
	// This is just a placeholder implementation

	query := `DELETE FROM reviews WHERE id = $1;`

	_, err := pool.Exec(ctx, query, r.ID)
	if err != nil {
		return err
	}

	return nil
}
