// ===============================================
// Test Suite: auth_service_test.go
// Description: Unit tests for authentication service business logic
//
// Test Groups:
//   - Setup and Test Data
//   - Login Success Tests
//   - Login Failure Tests
//   - User Management Tests
//   - Business Logic Validation Tests
// ===============================================

package auth

import (
	"context"
	"testing"

	"github.com/piyuo/go-counter-sample/internal/core/domain"
	"github.com/piyuo/go-counter-sample/internal/repository"
)

func TestAuthService_Login_Success(t *testing.T) {
	// Setup
	userRepo := repository.NewInMemoryUserRepository()
	authService := NewService(userRepo)
	ctx := context.Background()

	// Test successful login
	request := &domain.LoginRequest{
		Username: "user1",
		Password: "123",
	}

	response, err := authService.Login(ctx, request)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if !response.Success {
		t.Errorf("Expected success=true, got %v", response.Success)
	}

	if response.Message != "Login successful" {
		t.Errorf("Expected 'Login successful', got %v", response.Message)
	}

	if response.UserID != "1" {
		t.Errorf("Expected UserID='1', got %v", response.UserID)
	}

	if response.Token == "" {
		t.Error("Expected token to be generated")
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	// Setup
	userRepo := repository.NewInMemoryUserRepository()
	authService := NewService(userRepo)
	ctx := context.Background()

	// Test invalid password
	request := &domain.LoginRequest{
		Username: "user1",
		Password: "wrong_password",
	}

	response, err := authService.Login(ctx, request)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if response.Success {
		t.Errorf("Expected success=false, got %v", response.Success)
	}

	if response.Message != "Invalid username or password" {
		t.Errorf("Expected 'Invalid username or password', got %v", response.Message)
	}

	if response.UserID != "" {
		t.Errorf("Expected empty UserID, got %v", response.UserID)
	}

	if response.Token != "" {
		t.Errorf("Expected empty token, got %v", response.Token)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	// Setup
	userRepo := repository.NewInMemoryUserRepository()
	authService := NewService(userRepo)
	ctx := context.Background()

	// Test non-existent user
	request := &domain.LoginRequest{
		Username: "nonexistent",
		Password: "123",
	}

	response, err := authService.Login(ctx, request)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if response.Success {
		t.Errorf("Expected success=false, got %v", response.Success)
	}

	if response.Message != "Invalid username or password" {
		t.Errorf("Expected 'Invalid username or password', got %v", response.Message)
	}
}

func TestAuthService_Login_EmptyCredentials(t *testing.T) {
	// Setup
	userRepo := repository.NewInMemoryUserRepository()
	authService := NewService(userRepo)
	ctx := context.Background()

	// Test empty username
	request := &domain.LoginRequest{
		Username: "",
		Password: "123",
	}

	response, err := authService.Login(ctx, request)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if response.Success {
		t.Errorf("Expected success=false, got %v", response.Success)
	}

	if response.Message != "Username and password are required" {
		t.Errorf("Expected 'Username and password are required', got %v", response.Message)
	}
}

func TestAuthService_CreateUser(t *testing.T) {
	// Setup
	userRepo := repository.NewInMemoryUserRepository()
	authService := NewService(userRepo)
	ctx := context.Background()

	// Test creating new user
	newUser := &domain.User{
		Username: "testuser",
		Password: "testpass",
		Email:    "test@example.com",
	}

	createdUser, err := authService.CreateUser(ctx, newUser)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if createdUser == nil {
		t.Fatal("Expected created user, got nil")
	}

	if createdUser.Username != "testuser" {
		t.Errorf("Expected username='testuser', got %v", createdUser.Username)
	}

	if createdUser.Email != "test@example.com" {
		t.Errorf("Expected email='test@example.com', got %v", createdUser.Email)
	}

	if !createdUser.Active {
		t.Error("Expected user to be active")
	}

	// Password should be hashed
	if createdUser.Password == "testpass" {
		t.Error("Expected password to be hashed")
	}
}

func TestAuthService_CreateUser_AlreadyExists(t *testing.T) {
	// Setup
	userRepo := repository.NewInMemoryUserRepository()
	authService := NewService(userRepo)
	ctx := context.Background()

	// Try to create user that already exists
	existingUser := &domain.User{
		Username: "user1", // This user already exists in the sample data
		Password: "newpass",
		Email:    "new@example.com",
	}

	createdUser, err := authService.CreateUser(ctx, existingUser)

	// Assertions
	if err != domain.ErrUserAlreadyExists {
		t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
	}

	if createdUser != nil {
		t.Errorf("Expected nil user, got %v", createdUser)
	}
}
