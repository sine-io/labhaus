package workflow

import "context"

// Repository defines the interface for workflow persistence (DIP)
type Repository interface {
	// Create saves a new workflow
	Create(ctx context.Context, workflow *Entity) error

	// FindByID retrieves a workflow by ID
	FindByID(ctx context.Context, id string) (*Entity, error)

	// FindByUserID retrieves all workflows for a user
	FindByUserID(ctx context.Context, userID string, filter Filter) ([]*Entity, error)

	// Update modifies an existing workflow
	Update(ctx context.Context, workflow *Entity) error

	// Delete removes a workflow by ID
	Delete(ctx context.Context, id string) error

	// UpdateState updates only the workflow state
	UpdateState(ctx context.Context, id string, state State) error
}

// Filter represents query filters for workflow search
type Filter struct {
	State  *State
	Limit  int
	Offset int
}
