package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:password@localhost:5432/coin_control?sslmode=disable"
	}

	var err error
	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := initTables(); err != nil {
		return fmt.Errorf("failed to initialize tables: %w", err)
	}

	return nil
}

func initTables() error {
	ctx := context.Background()

	// Create users table
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT now()
	);`

	// Create auth table
	authTable := `
	CREATE TABLE IF NOT EXISTS auth (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		nickname TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT now()
	);`

	// Execute SQL commands
	if _, err := DB.Exec(ctx, usersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	if _, err := DB.Exec(ctx, authTable); err != nil {
		return fmt.Errorf("failed to create auth table: %w", err)
	}

	return nil
}
