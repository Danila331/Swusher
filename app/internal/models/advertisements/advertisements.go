package advertisements

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Advertisement структура для представления объявления.
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
	Address       string   `json:"address" db:"address"`
	Geolocation_X float64  `json:"geolocation_x" db:"geolocation_x"`
	Geolocation_Y float64  `json:"geolocation_y" db:"geolocation_y"`
}

// AdvertisementInterface интерфейс для работы с объявлениями.
type AdvertisementInterface interface {
	Create(ctx context.Context, pool *pgxpool.Pool) error
	Update(ctx context.Context, pool *pgxpool.Pool) error
	ReadByID(ctx context.Context, pool *pgxpool.Pool) (*Advertisement, error)
	ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Advertisement, error)
	ReadAllByUserID(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Advertisement, error)
	Delete(ctx context.Context, pool *pgxpool.Pool) error
}

// Create создает новое объявление в базе данных.
func (a *Advertisement) Create(ctx context.Context, pool *pgxpool.Pool) error {
	query := `INSERT INTO sharehub_advertisements (user_id, title, description, rental_rules, cost_per_day, cost_per_week, cost_per_month, photo_paths, category, address, geolocation_x, geolocation_y)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id;`

	err := pool.QueryRow(ctx, query, a.UserID, a.Title, a.Description, a.RentalRules,
		a.CostPerday,
		a.CostPerWeek,
		a.CostPerMonth,
		a.PhotoPaths,
		a.Category,
		a.Address,
		a.Geolocation_X,
		a.Geolocation_Y).Scan(&a.ID)
	if err != nil {
		return err
	}

	return nil
}

// Update обновляет существующее объявление в базе данных.
func (a *Advertisement) Update(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to update an existing advertisement in the database

	query := `UPDATE sharehub_advertisements SET user_id = $1, title = $2, description = $3, rental_rules = $4, cost_per_day = $5, cost_per_week = $6, cost_per_month = $7, photo_paths = $8, category = $9, address = $10, geolocation_x = $11, geolocation_y = $12 WHERE id = $13;`

	_, err := pool.Exec(ctx, query, a.UserID, a.Title, a.Description, a.RentalRules,
		a.CostPerday,
		a.CostPerWeek,
		a.CostPerMonth,
		a.PhotoPaths,
		a.Category,
		a.Address,
		a.Geolocation_X,
		a.Geolocation_Y,
		a.ID)
	if err != nil {
		return err
	}

	return nil
}

// ReadByID предоставляет объявление по ID из базы данных.
func (a *Advertisement) ReadByID(ctx context.Context, pool *pgxpool.Pool) error {
	var advertisement Advertisement
	query := `SELECT user_id, title, description, rental_rules, cost_per_day, cost_per_week, cost_per_month, photo_paths, category, address, geolocation_x, geolocation_y FROM sharehub_advertisements WHERE id = $1;`

	err := pool.QueryRow(ctx, query, a.ID).Scan(&advertisement.UserID,
		&a.Title,
		&a.Description,
		&a.RentalRules,
		&a.CostPerday,
		&a.CostPerWeek,
		&a.CostPerMonth,
		&a.PhotoPaths,
		&a.Category,
		&a.Address,
		&a.Geolocation_X,
		&a.Geolocation_Y)
	if err != nil {
		return err
	}

	return nil
}

// ReadAll предоставляет все объявления из базы данных с пагинацией.
func (a *Advertisement) ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Advertisement, error) {
	var advertisements []Advertisement
	query := `SELECT id, user_id, title, description, rental_rules, cost_per_day, cost_per_week, cost_per_month, photo_paths, category, address, geolocation_x, geolocation_y FROM sharehub_advertisements LIMIT $1 OFFSET $2;`

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
			&advertisement.Address,
			&advertisement.Geolocation_X,
			&advertisement.Geolocation_Y)
		if err != nil {
			return nil, err
		}
		advertisements = append(advertisements, advertisement)
	}

	return advertisements, nil
}

func (a *Advertisement) ReadAllByUserID(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Advertisement, error) {
	var advertisements []Advertisement
	query := `SELECT id, user_id, title, description, rental_rules, cost_per_day, cost_per_week, cost_per_month, photo_paths, category, address, geolocation_x, geolocation_y FROM sharehub_advertisements WHERE user_id = $1 LIMIT $2 OFFSET $3;`

	rows, err := pool.Query(ctx, query, a.UserID, limit, offset)
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
			&advertisement.Address,
			&advertisement.Geolocation_X,
			&advertisement.Geolocation_Y)
		if err != nil {
			return nil, err
		}
		advertisements = append(advertisements, advertisement)
	}

	return advertisements, nil
}

// Delete удаляет объявление из базы данных.
func (a *Advertisement) Delete(ctx context.Context, pool *pgxpool.Pool) error {
	query := `DELETE FROM sharehub_advertisements WHERE id = $1;`

	_, err := pool.Exec(ctx, query, a.ID)
	if err != nil {
		return err
	}

	return nil
}
