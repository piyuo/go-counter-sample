// ===============================================
// Module: grpc_server.go
// Description: gRPC-style server implementation for authentication service
//
// Sections:
//   - Server Structure
//   - Constructor
//   - Login Method Implementation
//   - Health Check Implementation
//   - Setup Functions
// ===============================================

package handlers

import (
	"context"
	"errors"

	"github.com/piyuo/go-counter-sample/internal/core/domain"
)

// GRPCServer implements a gRPC-style AuthService
// Note: This is a simplified implementation for demonstration
// In production, you would use the actual gRPC generated code
type GRPCServer struct {
	authService domain.AuthService
}

// NewGRPCServer creates a new gRPC-style server
func NewGRPCServer(authService domain.AuthService) *GRPCServer {
	return &GRPCServer{
		authService: authService,
	}
}

// gRPC-style message types for demonstration
// In real implementation, these would be generated from protobuf

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserId  string `json:"user_id"`
	Token   string `json:"token"`
}

type HealthRequest struct{}

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

// Login handles gRPC-style login requests
func (s *GRPCServer) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Validate request
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	// Convert gRPC request to domain request
	domainRequest := &domain.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	// Call business service
	response, err := s.authService.Login(ctx, domainRequest)
	if err != nil {
		return nil, err
	}

	// Convert domain response to gRPC response
	grpcResponse := &LoginResponse{
		Success: response.Success,
		Message: response.Message,
		UserId:  response.UserID,
		Token:   response.Token,
	}

	return grpcResponse, nil
}

// Health handles gRPC-style health check requests
func (s *GRPCServer) Health(ctx context.Context, req *HealthRequest) (*HealthResponse, error) {
	return &HealthResponse{
		Status:  "healthy",
		Service: "go-counter-sample",
		Version: "1.0.0",
	}, nil
}

// SetupGRPCServer sets up and returns a gRPC-style server
func SetupGRPCServer(authService domain.AuthService) *GRPCServer {
	return NewGRPCServer(authService)
}
