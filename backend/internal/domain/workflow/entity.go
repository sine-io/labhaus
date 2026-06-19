package workflow

import (
	"errors"
	"time"
)

// Entity represents a workflow execution in the domain
type Entity struct {
	ID        string
	UserID    string
	StyleID   string
	State     State
	Config    Config
	Result    *Result
	CreatedAt time.Time
	UpdatedAt time.Time
}

// State represents workflow execution state
type State string

const (
	StateDraft     State = "DRAFT"
	StatePending   State = "PENDING"
	StateRunning   State = "RUNNING"
	StatePaused    State = "PAUSED"
	StateCompleted State = "COMPLETED"
	StateFailed    State = "FAILED"
	StateCancelled State = "CANCELLED"
)

// Config holds workflow configuration
type Config struct {
	ImageCount int
	Width      int
	Height     int
	Steps      int
	Seed       int64
}

// Result holds workflow execution result
type Result struct {
	ImageURLs []string
	Duration  time.Duration
	Error     string
}

// Validation errors
var (
	ErrEmptyUserID      = errors.New("user ID cannot be empty")
	ErrEmptyStyleID     = errors.New("style ID cannot be empty")
	ErrInvalidState     = errors.New("invalid workflow state")
	ErrInvalidImageCount = errors.New("image count must be between 1 and 10")
	ErrInvalidDimensions = errors.New("image dimensions must be positive")
)

// New creates a new Workflow entity with validation
func New(userID, styleID string, config Config) (*Entity, error) {
	workflow := &Entity{
		UserID:    userID,
		StyleID:   styleID,
		State:     StateDraft,
		Config:    config,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := workflow.Validate(); err != nil {
		return nil, err
	}

	return workflow, nil
}

// Validate checks if the workflow entity is valid
func (w *Entity) Validate() error {
	if w.UserID == "" {
		return ErrEmptyUserID
	}
	if w.StyleID == "" {
		return ErrEmptyStyleID
	}
	if !w.State.IsValid() {
		return ErrInvalidState
	}
	if w.Config.ImageCount < 1 || w.Config.ImageCount > 10 {
		return ErrInvalidImageCount
	}
	if w.Config.Width <= 0 || w.Config.Height <= 0 {
		return ErrInvalidDimensions
	}
	return nil
}

// IsValid checks if state is valid
func (s State) IsValid() bool {
	switch s {
	case StateDraft, StatePending, StateRunning, StatePaused, StateCompleted, StateFailed, StateCancelled:
		return true
	}
	return false
}

// CanTransitionTo checks if state transition is allowed
func (w *Entity) CanTransitionTo(newState State) bool {
	validTransitions := map[State][]State{
		StateDraft:     {StatePending, StateCancelled},
		StatePending:   {StateRunning, StateCancelled},
		StateRunning:   {StatePaused, StateCompleted, StateFailed},
		StatePaused:    {StateRunning, StateCancelled},
		StateCompleted: {},
		StateFailed:    {},
		StateCancelled: {},
	}

	allowed, ok := validTransitions[w.State]
	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == newState {
			return true
		}
	}
	return false
}

// TransitionTo changes workflow state with validation
func (w *Entity) TransitionTo(newState State) error {
	if !w.CanTransitionTo(newState) {
		return errors.New("invalid state transition")
	}
	w.State = newState
	w.UpdatedAt = time.Now()
	return nil
}

// SetResult sets the workflow execution result
func (w *Entity) SetResult(result *Result) {
	w.Result = result
	w.UpdatedAt = time.Now()
}

// IsFinal checks if workflow is in a final state
func (w *Entity) IsFinal() bool {
	return w.State == StateCompleted || w.State == StateFailed || w.State == StateCancelled
}
