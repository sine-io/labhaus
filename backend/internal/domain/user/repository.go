package user

import "context"

// Repository defines the interface for user persistence (DIP)
type Repository interface {
	// Create saves a new user
	Create(ctx context.Context, user *Entity) error

	// FindByID retrieves a user by ID
	FindByID(ctx context.Context, id string) (*Entity, error)

	// FindByEmail retrieves a user by email
	FindByEmail(ctx context.Context, email string) (*Entity, error)

	// Update modifies an existing user
	Update(ctx context.Context, user *Entity) error

	// Delete removes a user by ID
	Delete(ctx context.Context, id string) error

	// Exists checks if a user with the given email exists
	Exists(ctx context.Context, email string) (bool, error)
}
