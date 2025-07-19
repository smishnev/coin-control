package main

import (
	"coin-control/backend/auth"
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx         context.Context
	authService *auth.AuthService
}

// NewApp creates a new App application struct
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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Login
func (a *App) Login(req auth.LoginRequest) (*auth.LoginResponse, error) {
	return a.authService.Login(a.ctx, req)
}

// ValidateToken JWT токен
func (a *App) ValidateToken(token string) (*auth.Claims, error) {
	return a.authService.ValidateToken(token)
}

// CreateAuth
func (a *App) CreateAuth(req auth.CreateAuthRequest) (*auth.Auth, error) {
	return a.authService.CreateAuth(a.ctx, req)
}

// GetAuthByID
func (a *App) GetAuthByID(id string) (*auth.Auth, error) {
	return a.authService.GetAuthByID(a.ctx, id)
}

// GetAuthByUserID
func (a *App) GetAuthByUserID(userID string) (*auth.Auth, error) {
	return a.authService.GetAuthByUserID(a.ctx, userID)
}

// UpdateAuth
func (a *App) UpdateAuth(req auth.UpdateAuthRequest) (*auth.Auth, error) {
	return a.authService.UpdateAuth(a.ctx, req)
}

// DeleteAuth
func (a *App) DeleteAuth(id string) error {
	return a.authService.DeleteAuth(a.ctx, id)
}

// GetAllAuth
func (a *App) GetAllAuth() ([]*auth.Auth, error) {
	return a.authService.GetAllAuth(a.ctx)
}

// CreateUserWithAuth create user with auth
func (a *App) CreateUserWithAuth(req auth.CreateAuthRequest, firstName, lastName string) (*auth.Auth, error) {
	return a.authService.CreateUserWithAuth(a.ctx, req, firstName, lastName)
}

// UpdatePasswordByNickname update password by nickname
func (a *App) UpdatePasswordByNickname(req auth.UpdatePasswordRequest) error {
	return a.authService.UpdatePasswordByNickname(a.ctx, req)
}

// ForgotPassword update password by nickname
func (a *App) ForgotPassword(req auth.ForgotPasswordRequest) error {
	return a.authService.ForgotPassword(a.ctx, req)
}
