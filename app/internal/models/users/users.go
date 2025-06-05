package users

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// User struct represents a user in the system.
type User struct {
	ID         string    `json:"id" db:"id"`
	Nickname   string    `json:"nickname" db:"nickname"`
	Name       string    `json:"name" db:"name"`
	LastName   string    `json:"last_name" db:"last_name"`
	Fatherland string    `json:"fatherland" db:"fatherland"`
	PhotoPath  string    `json:"photo_path" db:"photo_path"`
	Adress     string    `json:"adress" db:"adress"`
	Email      string    `json:"email" db:"email"`
	Phone      string    `json:"phone" db:"phone"`
	Password   string    `json:"password" db:"password"`
	Role       string    `json:"role" db:"role"`
	IsVerified bool      `json:"is_verified" db:"is_verified"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// UserInterface defines the methods that a user model should implement.
type UserInterface interface {
	Create(ctx context.Context, pool *pgxpool.Pool) error
	Update(ctx context.Context, pool *pgxpool.Pool) error
	ReadByID(ctx context.Context, pool *pgxpool.Pool) error
	ReadByEmail(ctx context.Context, pool *pgxpool.Pool) error
	ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]User, error)
	Delete(ctx context.Context, pool *pgxpool.Pool) error
}

// Create creates a new user in the database.
func (u *User) Create(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to create a new user in the database
	// This is just a placeholder implementation
	// Проверка на существование пользователя с такой почтой
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM sharehub_users WHERE email = $1);`
	err := pool.QueryRow(ctx, checkQuery, u.Email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("пользователь с такой почтой уже существует")
	}

	query := `INSERT INTO sharehub_users (nickname, name,  last_name, fatherland, photo_path, email, phone, password, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`

	err = pool.QueryRow(ctx, query, u.Nickname, u.Name, u.LastName, u.Fatherland, u.PhotoPath, u.Email, u.Phone, u.Password, u.Role).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

// Update updates an existing user in the database.
func (u *User) Update(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to update an existing user in the database
	// This is just a placeholder implementation

	query := `UPDATE users SET nickname = $1, name = $2, last_name = $3, fatherland = $4, photo_path=$5, email = $6, phone = $7, password = $8, role = $9 WHERE id = $10;`

	_, err := pool.Exec(ctx, query, u.Nickname, u.Name, u.LastName, u.Fatherland, u.PhotoPath, u.Email, u.Phone, u.Password, u.Role, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// ReadByID retrieves a user by ID from the database.
func (u *User) ReadByID(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to read a user by ID from the database
	// This is just a placeholder implementation
	query := `SELECT id, nickname, name, last_name, fatherland, photo_path, email, phone, password, role FROM sharehub_users WHERE id = $1;`

	err := pool.QueryRow(ctx, query, u.ID).Scan(&u.ID, &u.Nickname, &u.Name, &u.LastName, &u.Fatherland, &u.PhotoPath, &u.Email, &u.Phone, &u.Password, &u.Role)
	if err != nil {
		return err
	}

	return nil
}

// ReadByEmail retrieves a user by email from the database.
func (u *User) ReadByEmail(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to read a user by email from the database
	// This is just a placeholder implementation
	query := `SELECT id, nickname, name, last_name, fatherland, photo_path, email, phone, password, role FROM sharehub_users WHERE email = $1;`

	err := pool.QueryRow(ctx, query, u.Email).Scan(&u.ID, &u.Nickname, &u.Name, &u.LastName, &u.Fatherland, &u.PhotoPath, &u.Email, &u.Phone, &u.Password, &u.Role)
	if err != nil {
		return err
	}

	return nil
}

// ReadAll retrieves all users from the database with pagination.
func (u *User) ReadAll(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]User, error) {
	// Implement the logic to read all users from the database with pagination
	// This is just a placeholder implementation
	var users []User
	query := `SELECT id, nickname, name, last_name, fatherland, photo_path, email, phone, password, role FROM sharehub_users LIMIT $1 OFFSET $2;`

	rows, err := pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Nickname, &user.Name, &user.LastName, &user.Fatherland, &u.PhotoPath, &user.Email, &user.Phone, &user.Password, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// Delete removes a user from the database.
func (u *User) Delete(ctx context.Context, pool *pgxpool.Pool) error {
	// Implement the logic to delete a user from the database
	// This is just a placeholder implementation

	query := `DELETE FROM sharehub_users WHERE id = $1;`

	_, err := pool.Exec(ctx, query, u.ID)
	if err != nil {
		return err
	}

	return nil
}
