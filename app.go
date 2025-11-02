package main

import (
	"coin-control/backend/auth"
	"coin-control/backend/bybit"
	"coin-control/backend/queue"
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App represents the main application structure
type App struct {
	ctx                context.Context
	authService        *auth.AuthService
	bybitService       *bybit.BybitService
	priceSubscriptions map[string]chan bybit.PriceData
	priceMutex         sync.RWMutex
	queue              *queue.Queue
}

// NewApp creates a new App application instance
func NewApp() *App {
	return &App{
		authService:        auth.NewAuthService(),
		bybitService:       bybit.NewBybitService(),
		priceSubscriptions: make(map[string]chan bybit.PriceData),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// =============================================================================
// Utility methods
// =============================================================================

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// =============================================================================
// Authentication methods
// =============================================================================

// Login authenticates a user with email/password
func (a *App) Login(req auth.LoginRequest) (*auth.LoginResponse, error) {
	return a.authService.Login(a.ctx, req)
}

// ValidateToken validates a JWT token and returns claims
func (a *App) ValidateToken(token string) (*auth.Claims, error) {
	return a.authService.ValidateToken(token)
}

// =============================================================================
// User management methods
// =============================================================================

// CreateAuth creates a new authentication record
func (a *App) CreateAuth(req auth.CreateAuthRequest) (*auth.Auth, error) {
	return a.authService.CreateAuth(a.ctx, req)
}

// GetAuthByID retrieves authentication record by ID
func (a *App) GetAuthByID(id string) (*auth.Auth, error) {
	return a.authService.GetAuthByID(a.ctx, id)
}

// GetAuthByUserID retrieves authentication record by user ID
func (a *App) GetAuthByUserID(userID string) (*auth.Auth, error) {
	return a.authService.GetAuthByUserID(a.ctx, userID)
}

// UpdateAuth updates an existing authentication record
func (a *App) UpdateAuth(req auth.UpdateAuthRequest) (*auth.Auth, error) {
	return a.authService.UpdateAuth(a.ctx, req)
}

// DeleteAuth deletes an authentication record by ID
func (a *App) DeleteAuth(id string) error {
	return a.authService.DeleteAuth(a.ctx, id)
}

// GetAllAuth retrieves all authentication records
func (a *App) GetAllAuth() ([]*auth.Auth, error) {
	return a.authService.GetAllAuth(a.ctx)
}

// =============================================================================
// Advanced user operations
// =============================================================================

// CreateUserWithAuth creates a new user with authentication credentials
func (a *App) CreateUserWithAuth(req auth.CreateAuthRequest, firstName, lastName string) (*auth.Auth, error) {
	return a.authService.CreateUserWithAuth(a.ctx, req, firstName, lastName)
}

// UpdatePasswordByNickname updates user password using nickname
func (a *App) UpdatePasswordByNickname(req auth.UpdatePasswordRequest) error {
	return a.authService.UpdatePasswordByNickname(a.ctx, req)
}

// ForgotPassword handles password recovery process
func (a *App) ForgotPassword(req auth.ForgotPasswordRequest) error {
	return a.authService.ForgotPassword(a.ctx, req)
}

// =============================================================================
// Bybit integration methods
// =============================================================================

// FetchSpotHoldings fetches spot holdings for a user
func (a *App) FetchSpotHoldings(userId string) ([]bybit.Holding, error) {
	return a.bybitService.FetchSpotHoldings(userId)
}

// GetAssetBalance retrieves balance for a specific coin for the user
func (a *App) GetAssetBalance(userID string, coin string) (*bybit.CoinBalance, error) {
	return a.bybitService.GetAssetBalance(userID, coin)
}

// GetCoinIconURLs gets coin icon URLs
func (a *App) GetCoinIconURLs(coins []string) ([]bybit.IconEntry, error) {
	return a.bybitService.GetCoinIconURLs(coins)
}

// PrefetchCoinIcons prefetches coin icons for caching
func (a *App) PrefetchCoinIcons(coins []string) {
	a.bybitService.PrefetchCoinIcons(coins)
}

// =============================================================================
// Price streaming methods
// =============================================================================

// StartPriceStream starts streaming prices for a symbol
func (a *App) StartPriceStream(symbol string) error {
	a.priceMutex.Lock()
	defer a.priceMutex.Unlock()

	// Check if already subscribed
	if _, exists := a.priceSubscriptions[symbol]; exists {
		return nil
	}

	// Subscribe to price updates
	priceChan, err := a.bybitService.SubscribeToPrice(symbol)
	if err != nil {
		return err
	}

	// Store the channel
	a.priceSubscriptions[symbol] = priceChan

	// Start goroutine to handle price updates for this symbol
	go a.handlePriceUpdates(symbol, priceChan)

	return nil
}

// StopPriceStream stops streaming prices for a symbol
func (a *App) StopPriceStream(symbol string) {
	a.priceMutex.Lock()
	defer a.priceMutex.Unlock()

	if priceChan, exists := a.priceSubscriptions[symbol]; exists {
		a.bybitService.UnsubscribeFromPrice(symbol, priceChan)
		delete(a.priceSubscriptions, symbol)
	}
}

// GetCurrentPrice gets the current price for a symbol (one-time request)
func (a *App) GetCurrentPrice(symbol string) (string, error) {
	return a.bybitService.GetCurrentPrice(symbol)
}

// handlePriceUpdates processes incoming price updates and emits them to frontend
func (a *App) handlePriceUpdates(symbol string, priceChan chan bybit.PriceData) {
	for priceUpdate := range priceChan {
		// Use the original symbol parameter to ensure consistency
		// between frontend listener and backend emitter
		coinSymbol := strings.ToLower(symbol)

		priceValue := priceUpdate.Price
		if priceValue == "" {
			continue
		}

		// Emit event to frontend using the original symbol parameter
		eventName := fmt.Sprintf("price-update-%s", coinSymbol)
		eventData := map[string]interface{}{
			"symbol": coinSymbol,
			"price":  priceValue,
		}

		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, eventName, eventData)
		}
	}
}
