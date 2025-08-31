package user

import (
	"context"

	"coin-control/backend/database"

	"github.com/jackc/pgx/v5"
)

// =============================================================================
// Data structures
// =============================================================================

// User represents a user in the system
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// =============================================================================
// Service structure
// =============================================================================

// UserService provides user management functionality
type UserService struct{}

// NewUserService creates a new instance of UserService
func NewUserService() *UserService {
	return &UserService{}
}

// =============================================================================
// User operations
// =============================================================================

// CreateUserInTransaction creates a user within a database transaction
func (u *UserService) CreateUserInTransaction(ctx context.Context, tx pgx.Tx, user User) (string, error) {
	query := `
		INSERT INTO users (first_name, last_name)
		VALUES ($1, $2)
		RETURNING id
	`
	var newID string
	err := tx.QueryRow(ctx, query, user.FirstName, user.LastName).Scan(&newID)
	if err != nil {
		return "", err
	}
	return newID, nil
}

// UpdateUser creates or updates a user record
func (u *UserService) UpdateUser(user User) (string, error) {
	ctx := context.Background()

	query := `
		INSERT INTO users (id, first_name, last_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name
	`
	_, err := database.DB.Exec(ctx, query, user.ID, user.FirstName, user.LastName)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

// GetUser retrieves a user by ID
func (u *UserService) GetUser(id string) (*User, error) {
	ctx := context.Background()
	query := `
        SELECT id, first_name, last_name
        FROM users
        WHERE id = $1
    `
	var user User
	err := database.DB.QueryRow(ctx, query, id).Scan(&user.ID, &user.FirstName, &user.LastName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
