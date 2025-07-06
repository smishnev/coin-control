package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"coin-control/backend/database"
	"coin-control/backend/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Auth struct {
	ID           uuid.UUID `json:"id"`
	Nickname     string    `json:"nickname"`
	PasswordHash string    `json:"-"`
	UserID       uuid.UUID `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateAuthRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	UserID   string `json:"user_id"`
}

type UpdateAuthRequest struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Auth  *Auth  `json:"auth"`
	Token string `json:"token"`
}

type Claims struct {
	UserID   string `json:"user_id"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}

type AuthService struct {
	jwtSecret []byte
}

type UpdatePasswordRequest struct {
	Nickname    string `json:"nickname"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
type ForgotPasswordRequest struct {
	Nickname    string `json:"nickname"`
	NewPassword string `json:"new_password"`
}

func NewAuthService() *AuthService {
	secret := []byte("just-my-very-secret-key-RR-PP-OO")
	return &AuthService{
		jwtSecret: secret,
	}
}

// hashPassword
func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := sha256.Sum256(append([]byte(password), salt...))
	return hex.EncodeToString(hash[:]) + ":" + hex.EncodeToString(salt), nil
}

// verifyPassword
func verifyPassword(password, hashWithSalt string) bool {
	parts := strings.Split(hashWithSalt, ":")
	if len(parts) != 2 {
		return false
	}

	hash, saltHex := parts[0], parts[1]
	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return false
	}

	expectedHash := sha256.Sum256(append([]byte(password), salt...))
	return hex.EncodeToString(expectedHash[:]) == hash
}

// generateToken
func (s *AuthService) generateToken(auth *Auth) (string, error) {
	claims := &Claims{
		UserID:   auth.UserID.String(),
		Nickname: auth.Nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateToken
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// CreateAuth
func (s *AuthService) CreateAuth(ctx context.Context, req CreateAuthRequest) (*Auth, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	auth := &Auth{
		ID:           uuid.New(),
		Nickname:     req.Nickname,
		PasswordHash: passwordHash,
		UserID:       userID,
		CreatedAt:    time.Now(),
	}

	query := `
		INSERT INTO auth (id, nickname, password_hash, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, nickname, user_id, created_at
	`

	err = database.DB.QueryRow(ctx, query,
		auth.ID, auth.Nickname, auth.PasswordHash, auth.UserID, auth.CreatedAt,
	).Scan(&auth.ID, &auth.Nickname, &auth.UserID, &auth.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create auth: %w", err)
	}

	return auth, nil
}

// GetAuthByID
func (s *AuthService) GetAuthByID(ctx context.Context, id string) (*Auth, error) {
	authID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	query := `
		SELECT id, nickname, password_hash, user_id, created_at
		FROM auth
		WHERE id = $1
	`

	auth := &Auth{}
	err = database.DB.QueryRow(ctx, query, authID).Scan(
		&auth.ID, &auth.Nickname, &auth.PasswordHash, &auth.UserID, &auth.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("auth not found")
		}
		return nil, fmt.Errorf("failed to get auth: %w", err)
	}

	return auth, nil
}

// GetAuthByNickname
func (s *AuthService) GetAuthByNickname(ctx context.Context, nickname string) (*Auth, error) {
	query := `
		SELECT id, nickname, password_hash, user_id, created_at
		FROM auth
		WHERE nickname = $1
	`

	auth := &Auth{}
	err := database.DB.QueryRow(ctx, query, nickname).Scan(
		&auth.ID, &auth.Nickname, &auth.PasswordHash, &auth.UserID, &auth.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("auth not found")
		}
		return nil, fmt.Errorf("failed to get auth: %w", err)
	}

	return auth, nil
}

// UpdateAuth
func (s *AuthService) UpdateAuth(ctx context.Context, req UpdateAuthRequest) (*Auth, error) {
	authID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	// Get current record
	currentAuth, err := s.GetAuthByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Nickname != "" {
		currentAuth.Nickname = req.Nickname
	}
	if req.Password != "" {
		passwordHash, err := hashPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		currentAuth.PasswordHash = passwordHash
	}

	query := `
		UPDATE auth
		SET nickname = $1, password_hash = $2
		WHERE id = $3
		RETURNING id, nickname, user_id, created_at
	`

	err = database.DB.QueryRow(ctx, query,
		currentAuth.Nickname, currentAuth.PasswordHash, authID,
	).Scan(&currentAuth.ID, &currentAuth.Nickname, &currentAuth.UserID, &currentAuth.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update auth: %w", err)
	}

	return currentAuth, nil
}

// DeleteAuth
func (s *AuthService) DeleteAuth(ctx context.Context, id string) error {
	authID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	query := `DELETE FROM auth WHERE id = $1`
	result, err := database.DB.Exec(ctx, query, authID)
	if err != nil {
		return fmt.Errorf("failed to delete auth: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("auth not found")
	}

	return nil
}

// Login
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	auth, err := s.GetAuthByNickname(ctx, req.Nickname)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !verifyPassword(req.Password, auth.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate token
	token, err := s.generateToken(auth)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Auth:  auth,
		Token: token,
	}, nil
}

// GetAllAuth
func (s *AuthService) GetAllAuth(ctx context.Context) ([]*Auth, error) {
	query := `
		SELECT id, nickname, password_hash, user_id, created_at
		FROM auth
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth records: %w", err)
	}
	defer rows.Close()

	var auths []*Auth
	for rows.Next() {
		auth := &Auth{}
		err := rows.Scan(&auth.ID, &auth.Nickname, &auth.PasswordHash, &auth.UserID, &auth.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan auth record: %w", err)
		}
		auths = append(auths, auth)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating auth records: %w", err)
	}

	return auths, nil
}

// CreateAuthInTransaction create auth in transaction
func (s *AuthService) CreateAuthInTransaction(ctx context.Context, tx pgx.Tx, req CreateAuthRequest) (*Auth, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	auth := &Auth{
		ID:           uuid.New(),
		Nickname:     req.Nickname,
		PasswordHash: passwordHash,
		UserID:       userID,
		CreatedAt:    time.Now(),
	}

	query := `
		INSERT INTO auth (id, nickname, password_hash, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, nickname, user_id, created_at
	`

	err = tx.QueryRow(ctx, query,
		auth.ID, auth.Nickname, auth.PasswordHash, auth.UserID, auth.CreatedAt,
	).Scan(&auth.ID, &auth.Nickname, &auth.UserID, &auth.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create auth: %w", err)
	}

	return auth, nil
}

// CreateUserWithAuth create user with auth
func (s *AuthService) CreateUserWithAuth(ctx context.Context, req CreateAuthRequest, firstName, lastName string) (*Auth, error) {
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Rollback transaction in case of error

	// Create user in transaction
	userService := user.NewUserService()
	user := user.User{
		FirstName: firstName,
		LastName:  lastName,
	}

	userID, err := userService.CreateUserInTransaction(ctx, tx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Update user_id in request
	req.UserID = userID

	// Create auth record in transaction
	auth, err := s.CreateAuthInTransaction(ctx, tx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return auth, nil
}

// UpdatePasswordByNickname update password by nickname
func (s *AuthService) UpdatePasswordByNickname(ctx context.Context, req UpdatePasswordRequest) error {
	// Get auth record by nickname
	auth, err := s.GetAuthByNickname(ctx, req.Nickname)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Check old password
	if !verifyPassword(req.OldPassword, auth.PasswordHash) {
		return fmt.Errorf("invalid old password")
	}

	// Hash new password
	newPasswordHash, err := hashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	query := `
		UPDATE auth
		SET password_hash = $1
		WHERE nickname = $2
	`

	result, err := database.DB.Exec(ctx, query, newPasswordHash, req.Nickname)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// ForgotPasswordByNickname update password by nickname
func (s *AuthService) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {
	// Hash new password
	newPasswordHash, err := hashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	query := `
		UPDATE auth
		SET password_hash = $1
		WHERE nickname = $2
	`

	result, err := database.DB.Exec(ctx, query, newPasswordHash, req.Nickname)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
