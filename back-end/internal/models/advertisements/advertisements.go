package advertisements

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Advertisement struct represents an advertisement in the system.
type Advertisement struct {
	ID            string   `json:"id" db:"id"`
	UserID        string   `json:"user_id" db:"user_id"`
	Title         string   `json:"title" db:"title"`
	Description   string   `json:"description" db:"description"`
	RentalRules   string   `json:"rental_rules" db:"rental_rules"`
	CostPerday    float64  `json:"cost_per_day" db:"cost_per_day"`
	CostPerWeek   float64  `json:"cost_per_week" db:"cost_per_week"`
	CostPerMonth  float64  `json:"cost_per_month" db:"cost_per_month"`
	PhotoPaths    []string `json:"photo_paths" db:"photo_paths"`
	Category      string   `json:"category" db:"category"`
	Geolocation_X float64  `json:"geolocation_x" db:"geolocation_x"`
	Geolocation_Y float64  `json:"geolocation_y" db:"geolocation_y"`
}

// AdvertisementInterface defines the methods that a advertisement model should implement.
type AdvertisementInterface interface {
	Create(ctx context.Context, pool *pgxpool.Pool) error
	Update(ctx context.Context, pool *pgxpool.Pool) error
	ReadByID(ctx context.Context, pool *pgxpool.Pool) (*Advertisement, error)
	ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Advertisement, error)
	Delete(ctx context.Context, pool *pgxpool.Pool) error
}

// Create creates a new advertisement in the database.
func (a *Advertisement) Create(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to create a new advertisement in the database
	// This is just a placeholder implementation

	query := `INSERT INTO advertisements (user_id, title, description, rental_rules, cost_per_day, cost_per_week, cost_per_month, photo_paths, category, geolocation_x, geolocation_y)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;`

	err := pool.QueryRow(ctx, query, a.UserID, a.Title, a.Description, a.RentalRules,
		a.CostPerday,
		a.CostPerWeek,
		a.CostPerMonth,
		a.PhotoPaths,
		a.Category,
		a.Geolocation_X,
		a.Geolocation_Y).Scan(&a.ID)
	if err != nil {
		return err
	}

	return nil
}

// Update updates an existing advertisement in the database.
func (a *Advertisement) Update(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to update an existing advertisement in the database
	// This is just a placeholder implementation

	query := `UPDATE advertisements SET user_id = $1, title = $2, description = $3, rental_rules = $4, cost_per_day = $5, cost_per_week = $6, cost_per_month = $7, photo_paths = $8, category = $9, geolocation_x = $10, geolocation_y = $11 WHERE id = $12;`

	_, err := pool.Exec(ctx, query, a.UserID, a.Title, a.Description, a.RentalRules,
		a.CostPerday,
		a.CostPerWeek,
		a.CostPerMonth,
		a.PhotoPaths,
		a.Category,
		a.Geolocation_X,
		a.Geolocation_Y,
		a.ID)
	if err != nil {
		return err
	}

	return nil
}

// ReadByID retrieves an advertisement by its ID from the database.
func (a *Advertisement) ReadByID(ctx context.Context, pool *pgxpool.Pool) (*Advertisement, error) {
	// Implement the logic to retrieve an advertisement by its ID from the database
	// This is just a placeholder implementation
	var advertisement Advertisement
	query := `SELECT user_id, title, description, rental_rules, cost_per_day, cost_per_week, cost_per_month, photo_paths, category, geolocation_x, geolocation_y FROM advertisements WHERE id = $1;`

	err := pool.QueryRow(ctx, query, a.ID).Scan(&advertisement.UserID,
		&advertisement.Title,
		&advertisement.Description,
		&advertisement.RentalRules,
		&advertisement.CostPerday,
		&advertisement.CostPerWeek,
		&advertisement.CostPerMonth,
		&advertisement.PhotoPaths,
		&advertisement.Category,
		&advertisement.Geolocation_X,
		&advertisement.Geolocation_Y)
	if err != nil {
		return nil, err
	}

	return &advertisement, nil
}

// ReadAll retrieves all advertisements from the database with pagination.
func (a *Advertisement) ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Advertisement, error) {
	// Implement the logic to retrieve all advertisements from the database with pagination
	// This is just a placeholder implementation
	var advertisements []Advertisement
	query := `SELECT id, user_id, title, description, rental_rules, cost_per_day, cost_per_week, cost_per_month, photo_paths, category, geolocation_x, geolocation_y FROM advertisements LIMIT $1 OFFSET $2;`

	rows, err := pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var advertisement Advertisement
		err := rows.Scan(&advertisement.ID,
			&advertisement.UserID,
			&advertisement.Title,
			&advertisement.Description,
			&advertisement.RentalRules,
			&advertisement.CostPerday,
			&advertisement.CostPerWeek,
			&advertisement.CostPerMonth,
			&advertisement.PhotoPaths,
			&advertisement.Category,
			&advertisement.Geolocation_X,
			&advertisement.Geolocation_Y)
		if err != nil {
			return nil, err
		}
		advertisements = append(advertisements, advertisement)
	}

	return advertisements, nil
}

// Delete deletes an advertisement from the database.
func (a *Advertisement) Delete(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to delete an advertisement from the database
	// This is just a placeholder implementation

	query := `DELETE FROM advertisements WHERE id = $1;`

	_, err := pool.Exec(ctx, query, a.ID)
	if err != nil {
		return err
	}

	return nil
}
