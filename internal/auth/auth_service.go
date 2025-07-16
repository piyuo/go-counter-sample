// ===============================================
// Module: auth_service.go
// Description: Authentication business service implementation
//
// Sections:
//   - Service Structure
//   - Constructor
//   - Login Business Logic
//   - User Validation
//   - User Management Methods
//   - Helper Functions
// ===============================================

package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/piyuo/go-counter-sample/internal/core/domain"
)

// Service implements the AuthService interface
type Service struct {
	userRepo domain.UserRepository
}

// NewService creates a new authentication service
func NewService(userRepo domain.UserRepository) domain.AuthService {
	return &Service{
		userRepo: userRepo,
	}
}

// Login handles user authentication
func (s *Service) Login(ctx context.Context, request *domain.LoginRequest) (*domain.LoginResponse, error) {
	// Validate input
	if request.Username == "" || request.Password == "" {
		return &domain.LoginResponse{
			Success: false,
			Message: "Username and password are required",
		}, nil
	}

	// Validate user credentials
	user, err := s.ValidateUser(ctx, request.Username, request.Password)
	if err != nil {
		if err == domain.ErrUserNotFound || err == domain.ErrInvalidCredentials {
			return &domain.LoginResponse{
				Success: false,
				Message: "Invalid username or password",
			}, nil
		}
		return nil, fmt.Errorf("login failed: %w", err)
	}

	// Check if user is active
	if !user.Active {
		return &domain.LoginResponse{
			Success: false,
			Message: "Account is inactive",
		}, nil
	}

	// Successful login
	return &domain.LoginResponse{
		Success: true,
		Message: "Login successful",
		UserID:  user.ID,
		Token:   s.generateToken(user), // In real app, use JWT
	}, nil
}

// ValidateUser validates user credentials
func (s *Service) ValidateUser(ctx context.Context, username, password string) (*domain.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// In a real application, you would hash the password and compare
	// For demo purposes, we're doing simple comparison
	if user.Password != s.hashPassword(password) {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}

// CreateUser creates a new user
func (s *Service) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Hash password before storing
	user.Password = s.hashPassword(user.Password)
	user.Active = true

	return s.userRepo.Create(ctx, user)
}

// GetUserByUsername retrieves a user by username
func (s *Service) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

// Helper functions

// hashPassword hashes a password (simplified for demo)
func (s *Service) hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// generateToken generates a simple token (in real app, use JWT)
func (s *Service) generateToken(user *domain.User) string {
	return fmt.Sprintf("token_%s_%s", user.ID, strings.ToLower(user.Username))
}
