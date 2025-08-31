package bybit

import (
	"coin-control/backend/database"
	"context"
	"time"
)

type Bybit struct {
	ID        string    `json:"id"`
	ApiKey    string    `json:"apiKey"`
	ApiSecret string    `json:"apiSecret"`
	UserId    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type BybitService struct {
}

func NewBybitService() *BybitService {
	return &BybitService{}
}

// Expose streaming & holdings to Wails using distinct method names to avoid recursion
func (s *BybitService) FetchSpotHoldings(userId string) ([]Holding, error) {
	return s.getSpotHoldings(context.Background(), userId)
}

// Expose coin icons with cache
func (s *BybitService) GetCoinIcons(coins []string) ([]IconEntry, error) {
	return s.getCoinIcons(coins)
}

// Fast, URL-only variant (no data URLs)
func (s *BybitService) GetCoinIconURLs(coins []string) ([]IconEntry, error) {
	return s.getCoinIconURLs(coins)
}

// Background prefetch to warm disk cache for icons
func (s *BybitService) PrefetchCoinIcons(coins []string) {
	go s.prefetchCoinIcons(coins)
}

func (s *BybitService) CreateBybit(bybit Bybit) (string, error) {
	ctx := context.Background()

	// encrypt secret on write
	encSecret, err := encryptString(bybit.ApiSecret)
	if err != nil {
		return "", err
	}

	query := `
		INSERT INTO bybit (api_key, api_secret, user_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var newID string
	err = database.DB.QueryRow(ctx, query, bybit.ApiKey, encSecret, bybit.UserId).Scan(&newID)
	if err != nil {
		return "", err
	}
	return newID, nil
}

func (s *BybitService) UpsertBybit(bybitApi string, bybitApiSecret string, userId string) error {
	ctx := context.Background()

	query := `
		INSERT INTO bybit (api_key, api_secret, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE SET
			api_key = EXCLUDED.api_key,
			api_secret = EXCLUDED.api_secret,
			updated_at = EXCLUDED.updated_at
	`
	now := time.Now()
	// encrypt secret on upsert
	encSecret, err := encryptString(bybitApiSecret)
	if err != nil {
		return err
	}
	_, err = database.DB.Exec(ctx, query, bybitApi, encSecret, userId, now, now)
	if err != nil {
		return err
	}
	return nil
}

func (s *BybitService) GetBybitByUserId(userId string) (*Bybit, error) {
	ctx := context.Background()

	query := `
		SELECT id, api_key, api_secret, user_id, created_at, updated_at
		FROM bybit
		WHERE user_id = $1
	`
	var bybit Bybit
	err := database.DB.QueryRow(ctx, query, userId).Scan(&bybit.ID, &bybit.ApiKey, &bybit.ApiSecret, &bybit.UserId, &bybit.CreatedAt, &bybit.UpdatedAt)
	if err != nil {
		return nil, err
	}
	// decrypt on read; if not encrypted yet, passthrough happens
	if dec, derr := decryptString(bybit.ApiSecret); derr == nil {
		bybit.ApiSecret = dec
	}
	return &bybit, nil
}
