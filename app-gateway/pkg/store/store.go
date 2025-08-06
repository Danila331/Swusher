package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DBConfig holds the configuration for the database connection.
type DBConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	ConnString      string
	MaxConns        int
	MinConns        int
	MaxConnIdleTime time.Duration
	MaxConnLifetime time.Duration
}

// NewDBConfig creates a new DBConfig with default values.
func NewDBConfig(Host string,
	Port int,
	User,
	Password,
	DBName,
	SSLMode string,
	ConnString string,
	MaxConns,
	MinConns int,
	MaxConnIdleTime,
	MaxConnLifetime time.Duration) DBConfig {
	return DBConfig{
		Host:            Host,
		Port:            Port,
		User:            User,
		Password:        Password,
		DBName:          DBName,
		SSLMode:         SSLMode,
		ConnString:      ConnString,
		MaxConns:        MaxConns,
		MinConns:        MinConns,
		MaxConnIdleTime: MaxConnIdleTime,
		MaxConnLifetime: MaxConnLifetime,
	}
}

// NewPgxPool creates a new DBConfig with default values.
func NexPgxPool(ctx context.Context, cfg DBConfig) (*pgxpool.Pool, error) {
	// Create a new pgxpool.Config
	poolConfig, err := pgxpool.ParseConfig(cfg.ConnString)
	if err != nil {
		return nil, err
	}

	// Set the pool configuration options
	poolConfig.MaxConns = int32(cfg.MaxConns)
	poolConfig.MinConns = int32(cfg.MinConns)
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime

	// Create a new pgxpool.Pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

// CreateTables creates the necessary tables in the database.
func CreateTables(ctx context.Context, pool *pgxpool.Pool) error {

	// Create the users table
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS sharehub_users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		nickname    VARCHAR(50) UNIQUE DEFAULT NULL,
		name        VARCHAR(100),
		last_name   VARCHAR(100),
		fatherland  VARCHAR(100),
		photo_path  TEXT,
		adress      TEXT,
		email       VARCHAR(255) UNIQUE NOT NULL,
		phone       VARCHAR(20),
		password    TEXT NOT NULL,
		role        VARCHAR(50) NOT NULL DEFAULT 'user',
		is_verified BOOLEAN NOT NULL DEFAULT FALSE,
		created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`)

	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create the passports table
	_, err = pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS sharehub_passports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES sharehub_users(id) ON DELETE CASCADE,
    photo_path_with_author    TEXT,
    photo_path_without_author TEXT,
    status      VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`)

	if err != nil {
		return fmt.Errorf("failed to create passports table: %w", err)
	}

	// Create the advertisements table
	_, err = pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS sharehub_advertisements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES sharehub_users(id) ON DELETE CASCADE,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    rental_rules    TEXT,
    cost_per_day    DOUBLE PRECISION NOT NULL DEFAULT 0,
    cost_per_week   DOUBLE PRECISION NOT NULL DEFAULT 0,
    cost_per_month  DOUBLE PRECISION NOT NULL DEFAULT 0,
    photo_paths     TEXT[] NOT NULL DEFAULT '{}',
    category        VARCHAR(100),
    address         TEXT,
	geolocation_x   DOUBLE PRECISION,
    geolocation_y   DOUBLE PRECISION
);`)

	if err != nil {
		return fmt.Errorf("failed to create advertisements table: %w", err)
	}

	// Create the chats table
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS sharehub_chats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_from_id UUID NOT NULL REFERENCES sharehub_users(id) ON DELETE CASCADE,
    user_to_id   UUID NOT NULL REFERENCES sharehub_users(id) ON DELETE CASCADE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`)

	if err != nil {
		return fmt.Errorf("failed to create chats table: %w", err)
	}

	// Create the messages table
	_, err = pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS sharehub_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id UUID NOT NULL REFERENCES sharehub_chats(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES sharehub_users(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`)
	if err != nil {
		return fmt.Errorf("failed to create messages table: %w", err)
	}

	// Create the reviews table
	_, err = pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS sharehub_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    advertisement_id UUID NOT NULL REFERENCES sharehub_advertisements(id) ON DELETE CASCADE,
    user_from_id UUID NOT NULL REFERENCES sharehub_users(id) ON DELETE CASCADE,
    score NUMERIC(2,1) CHECK (score >= 0 AND score <= 5) NOT NULL,
    text TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`)

	if err != nil {
		return fmt.Errorf("failed to create reviews table: %w", err)
	}
	return nil
}
