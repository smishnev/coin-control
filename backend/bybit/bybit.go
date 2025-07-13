package bybit

import (
	"coin-control/backend/database"
	"context"
	"time"
)

type Bybit struct {
	ID        string    `json:"id"`
	ApiKey    string    `json:"apiKey"`
	UserId    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type BybitService struct {
}

func NewBybitService() *BybitService {
	return &BybitService{}
}

func (s *BybitService) CreateBybit(bybit Bybit) (string, error) {
	ctx := context.Background()

	query := `
		INSERT INTO bybit (api_key, user_id)
		VALUES ($1, $2)
		RETURNING id
	`
	var newID string
	err := database.DB.QueryRow(ctx, query, bybit.ApiKey, bybit.UserId).Scan(&newID)
	if err != nil {
		return "", err
	}
	return newID, nil
}

func (s *BybitService) UpsertBybit(bybitApi string, userId string) error {
	ctx := context.Background()

	query := `
		INSERT INTO bybit (api_key, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			api_key = EXCLUDED.api_key,
			updated_at = EXCLUDED.updated_at
	`
	now := time.Now()
	_, err := database.DB.Exec(ctx, query, bybitApi, userId, now, now)
	if err != nil {
		return err
	}
	return nil
}

func (s *BybitService) GetBybitByUserId(userId string) (*Bybit, error) {
	ctx := context.Background()

	query := `
		SELECT id, api_key, user_id, created_at, updated_at
		FROM bybit
		WHERE user_id = $1
	`
	var bybit Bybit
	err := database.DB.QueryRow(ctx, query, userId).Scan(&bybit.ID, &bybit.ApiKey, &bybit.UserId, &bybit.CreatedAt, &bybit.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &bybit, nil
}
