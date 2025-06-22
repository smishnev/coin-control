package user

import (
	"context"

	"coin-control/backend/database"
)

type User struct {
	ID        *string `json:"id,omitempty"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
}

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) CreateOrUpdate(user User) (string, error) {
	ctx := context.Background()

	if user.ID == nil || *user.ID == "" {
		query := `
			INSERT INTO users (first_name, last_name)
			VALUES ($1, $2)
			RETURNING id
		`
		var newID string
		err := database.DB.QueryRow(ctx, query, user.FirstName, user.LastName).Scan(&newID)
		if err != nil {
			return "", err
		}
		return newID, nil
	}

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
	return *user.ID, nil
}

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
