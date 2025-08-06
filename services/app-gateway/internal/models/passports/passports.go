package passports

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Passport структура представляет собой паспорт пользователя.
type Passport struct {
	ID                     string    `json:"id" db:"id"`
	UserID                 string    `json:"user_id" db:"user_id"`
	PhotoPathWithAuthor    string    `json:"photo_path_with_author" db:"photo_path_with_author"`
	PhotoPathWithoutAuthor string    `json:"photo_path_without_author" db:"photo_path_without_author"`
	Status                 string    `json:"status" db:"status"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

// PassportInterface определяет методы, которые должен реализовать модель паспорта.
type PassportInterface interface {
	Create(ctx context.Context, pool *pgxpool.Pool) error
	Update(ctx context.Context, pool *pgxpool.Pool) error
	ReadByID(ctx context.Context, pool *pgxpool.Pool) (*Passport, error)
	ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Passport, error)
	Delete(ctx context.Context, pool *pgxpool.Pool) error
}

// Create создает новый паспорт в базе данных.
func (p *Passport) Create(ctx context.Context, pool *pgxpool.Pool) error {
	query := `INSERT INTO passports (user_id, photo_path_with_author, photo_path_without_author, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

	err := pool.QueryRow(ctx, query, p.UserID, p.PhotoPathWithAuthor, p.PhotoPathWithoutAuthor, p.Status, p.CreatedAt, p.UpdatedAt).Scan(&p.ID)
	if err != nil {
		return err
	}

	return nil
}

// Update обновляет существующий паспорт в базе данных.
func (p *Passport) Update(ctx context.Context, pool *pgxpool.Pool) error {
	query := `UPDATE passports SET user_id = $1, photo_path_with_author = $2, photo_path_without_author = $3, status = $4, created_at = $5, updated_at = $6 WHERE id = $7;`

	_, err := pool.Exec(ctx, query, p.UserID, p.PhotoPathWithAuthor, p.PhotoPathWithoutAuthor, p.Status, p.CreatedAt, p.UpdatedAt, p.ID)
	if err != nil {
		return err
	}

	return nil
}

// ReadByID предоставляет паспорт по ID из базы данных.
func (p *Passport) ReadByID(ctx context.Context, pool *pgxpool.Pool) (*Passport, error) {
	query := `SELECT user_id, photo_path_with_author, photo_path_without_author, status, created_at, updated_at FROM passports WHERE id = $1;`

	err := pool.QueryRow(ctx, query, p.ID).Scan(&p.UserID, &p.PhotoPathWithAuthor, &p.PhotoPathWithoutAuthor, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// ReadAll предоставляет все паспорта из базы данных с пагинацией.
func (p *Passport) ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]Passport, error) {
	query := `SELECT id, user_id, photo_path_with_author, photo_path_without_author, status, created_at, updated_at FROM passports LIMIT $1 OFFSET $2;`

	rows, err := pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passports []Passport
	for rows.Next() {
		var passport Passport
		err := rows.Scan(&passport.ID, &passport.UserID, &passport.PhotoPathWithAuthor, &passport.PhotoPathWithoutAuthor, &passport.Status, &passport.CreatedAt, &passport.UpdatedAt)
		if err != nil {
			return nil, err
		}
		passports = append(passports, passport)
	}

	return passports, nil
}

// Delete удаляет паспорт из базы данных.
func (p *Passport) Delete(ctx context.Context, pool *pgxpool.Pool) error {
	query := `DELETE FROM passports WHERE id = $1;`

	_, err := pool.Exec(ctx, query, p.ID)
	if err != nil {
		return err
	}

	return nil
}
