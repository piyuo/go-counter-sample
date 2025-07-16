// ===============================================
// Module: user.go
// Description: User domain model and related types for authentication
//
// Sections:
//   - User Domain Model
//   - Login Request/Response Types
//   - Authentication Interfaces
//   - User Repository Interface
// ===============================================

package domain

import (
	"context"
	"errors"
)

// User represents a user in the system
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // In real app, this would be hashed
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}

// LoginRequest represents a login request payload
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  string `json:"user_id,omitempty"`
	Token   string `json:"token,omitempty"` // For JWT or session token
}

// AuthService defines the business logic interface for authentication
type AuthService interface {
	Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error)
	ValidateUser(ctx context.Context, username, password string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
}

// UserRepository defines the data access interface for users
type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id string) error
}

// Common domain errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserInactive       = errors.New("user account is inactive")
)
