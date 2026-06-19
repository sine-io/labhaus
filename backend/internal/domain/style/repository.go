package style

import "context"

// Repository defines the interface for style persistence (DIP - Dependency Inversion Principle)
// Domain layer defines the interface, infrastructure layer implements it
type Repository interface {
	// Create saves a new style
	Create(ctx context.Context, style *Entity) error

	// FindByID retrieves a style by ID
	FindByID(ctx context.Context, id string) (*Entity, error)

	// FindAll retrieves all styles with optional filters
	FindAll(ctx context.Context, filter Filter) ([]*Entity, error)

	// Update modifies an existing style
	Update(ctx context.Context, style *Entity) error

	// Delete removes a style by ID
	Delete(ctx context.Context, id string) error

	// Search performs full-text search on styles
	Search(ctx context.Context, query string, limit int) ([]*Entity, error)
}

// Filter represents query filters for style search
type Filter struct {
	Category string
	Tags     []string
	Limit    int
	Offset   int
}
