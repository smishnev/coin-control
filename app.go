package main

import (
	"coin-control/backend/auth"
	"context"
	"fmt"
)

// App represents the main application structure
type App struct {
	ctx         context.Context
	authService *auth.AuthService
}

// NewApp creates a new App application instance
func NewApp() *App {
	return &App{
		authService: auth.NewAuthService(),
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
