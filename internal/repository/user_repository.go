// ===============================================
// Module: user_repository.go
// Description: In-memory user repository implementation for demo purposes
//
// Sections:
//   - Repository Structure
//   - Constructor and Sample Data
//   - CRUD Operations
//   - Helper Methods
// ===============================================

package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/piyuo/go-counter-sample/internal/core/domain"
)

// InMemoryUserRepository implements UserRepository interface
type InMemoryUserRepository struct {
	users  map[string]*domain.User
	mutex  sync.RWMutex
	nextID int
}

// NewInMemoryUserRepository creates a new in-memory user repository with sample data
func NewInMemoryUserRepository() domain.UserRepository {
	repo := &InMemoryUserRepository{
		users:  make(map[string]*domain.User),
		mutex:  sync.RWMutex{},
		nextID: 1,
	}

	// Add sample user for testing (password is "123" hashed)
	// Hash of "123" using SHA256
	hashedPassword := "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3"

	sampleUser := &domain.User{
		ID:       "1",
		Username: "user1",
		Password: hashedPassword,
		Email:    "user1@example.com",
		Active:   true,
	}

	repo.users[sampleUser.Username] = sampleUser
	repo.nextID = 2

	return repo
}

// GetByUsername retrieves a user by username
func (r *InMemoryUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[username]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	// Return a copy to avoid external modifications
	userCopy := *user
	return &userCopy, nil
}

// Create creates a new user
func (r *InMemoryUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if user already exists
	if _, exists := r.users[user.Username]; exists {
		return nil, domain.ErrUserAlreadyExists
	}

	// Generate ID if not provided
	if user.ID == "" {
		user.ID = fmt.Sprintf("%d", r.nextID)
		r.nextID++
	}

	// Store user
	userCopy := *user
	r.users[user.Username] = &userCopy

	return &userCopy, nil
}

// Update updates an existing user
func (r *InMemoryUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[user.Username]; !exists {
		return nil, domain.ErrUserNotFound
	}

	// Update user
	userCopy := *user
	r.users[user.Username] = &userCopy

	return &userCopy, nil
}

// Delete deletes a user by ID
func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Find user by ID
	for username, user := range r.users {
		if user.ID == id {
			delete(r.users, username)
			return nil
		}
	}

	return domain.ErrUserNotFound
}
