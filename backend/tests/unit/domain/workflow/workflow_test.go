package workflow_test

import (
	"testing"
	"time"

	"github.com/labhaus/backend/internal/domain/workflow"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		styleID string
		config  workflow.Config
		wantErr error
	}{
		{
			name:    "valid workflow",
			userID:  "user-123",
			styleID: "style-456",
			config: workflow.Config{
				ImageCount: 4,
				Width:      512,
				Height:     512,
				Steps:      30,
				Seed:       12345,
			},
			wantErr: nil,
		},
		{
			name:    "empty user ID",
			userID:  "",
			styleID: "style-456",
			config: workflow.Config{
				ImageCount: 1,
				Width:      512,
				Height:     512,
			},
			wantErr: workflow.ErrEmptyUserID,
		},
		{
			name:    "empty style ID",
			userID:  "user-123",
			styleID: "",
			config: workflow.Config{
				ImageCount: 1,
				Width:      512,
				Height:     512,
			},
			wantErr: workflow.ErrEmptyStyleID,
		},
		{
			name:    "invalid image count - too low",
			userID:  "user-123",
			styleID: "style-456",
			config: workflow.Config{
				ImageCount: 0,
				Width:      512,
				Height:     512,
			},
			wantErr: workflow.ErrInvalidImageCount,
		},
		{
			name:    "invalid image count - too high",
			userID:  "user-123",
			styleID: "style-456",
			config: workflow.Config{
				ImageCount: 11,
				Width:      512,
				Height:     512,
			},
			wantErr: workflow.ErrInvalidImageCount,
		},
		{
			name:    "invalid dimensions",
			userID:  "user-123",
			styleID: "style-456",
			config: workflow.Config{
				ImageCount: 1,
				Width:      0,
				Height:     512,
			},
			wantErr: workflow.ErrInvalidDimensions,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := workflow.New(tt.userID, tt.styleID, tt.config)
			if err != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && w == nil {
				t.Error("New() returned nil workflow with no error")
			}
			if err == nil {
				if w.State != workflow.StateDraft {
					t.Errorf("Initial state = %v, want DRAFT", w.State)
				}
				if w.CreatedAt.IsZero() {
					t.Error("CreatedAt should not be zero")
				}
				if w.UpdatedAt.IsZero() {
					t.Error("UpdatedAt should not be zero")
				}
			}
		})
	}
}

func TestState_IsValid(t *testing.T) {
	validStates := []workflow.State{
		workflow.StateDraft,
		workflow.StatePending,
		workflow.StateRunning,
		workflow.StatePaused,
		workflow.StateCompleted,
		workflow.StateFailed,
		workflow.StateCancelled,
	}

	for _, state := range validStates {
		if !state.IsValid() {
			t.Errorf("State %v should be valid", state)
		}
	}

	invalidState := workflow.State("INVALID")
	if invalidState.IsValid() {
		t.Error("Invalid state should return false")
	}
}

func TestEntity_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name      string
		fromState workflow.State
		toState   workflow.State
		want      bool
	}{
		// Valid transitions
		{name: "DRAFT -> PENDING", fromState: workflow.StateDraft, toState: workflow.StatePending, want: true},
		{name: "DRAFT -> CANCELLED", fromState: workflow.StateDraft, toState: workflow.StateCancelled, want: true},
		{name: "PENDING -> RUNNING", fromState: workflow.StatePending, toState: workflow.StateRunning, want: true},
		{name: "PENDING -> CANCELLED", fromState: workflow.StatePending, toState: workflow.StateCancelled, want: true},
		{name: "RUNNING -> PAUSED", fromState: workflow.StateRunning, toState: workflow.StatePaused, want: true},
		{name: "RUNNING -> COMPLETED", fromState: workflow.StateRunning, toState: workflow.StateCompleted, want: true},
		{name: "RUNNING -> FAILED", fromState: workflow.StateRunning, toState: workflow.StateFailed, want: true},
		{name: "PAUSED -> RUNNING", fromState: workflow.StatePaused, toState: workflow.StateRunning, want: true},
		{name: "PAUSED -> CANCELLED", fromState: workflow.StatePaused, toState: workflow.StateCancelled, want: true},

		// Invalid transitions
		{name: "DRAFT -> RUNNING", fromState: workflow.StateDraft, toState: workflow.StateRunning, want: false},
		{name: "PENDING -> COMPLETED", fromState: workflow.StatePending, toState: workflow.StateCompleted, want: false},
		{name: "COMPLETED -> RUNNING", fromState: workflow.StateCompleted, toState: workflow.StateRunning, want: false},
		{name: "FAILED -> RUNNING", fromState: workflow.StateFailed, toState: workflow.StateRunning, want: false},
		{name: "CANCELLED -> RUNNING", fromState: workflow.StateCancelled, toState: workflow.StateRunning, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &workflow.Entity{
				UserID:  "user-123",
				StyleID: "style-456",
				State:   tt.fromState,
			}
			got := w.CanTransitionTo(tt.toState)
			if got != tt.want {
				t.Errorf("CanTransitionTo(%v) = %v, want %v", tt.toState, got, tt.want)
			}
		})
	}
}

func TestEntity_TransitionTo(t *testing.T) {
	w, _ := workflow.New("user-123", "style-456", workflow.Config{
		ImageCount: 1,
		Width:      512,
		Height:     512,
	})

	originalUpdatedAt := w.UpdatedAt
	time.Sleep(1 * time.Millisecond) // Ensure time difference

	// Valid transition
	err := w.TransitionTo(workflow.StatePending)
	if err != nil {
		t.Errorf("TransitionTo(PENDING) should succeed, got error: %v", err)
	}
	if w.State != workflow.StatePending {
		t.Errorf("State = %v, want PENDING", w.State)
	}
	if w.UpdatedAt == originalUpdatedAt {
		t.Error("UpdatedAt should be updated after state transition")
	}

	// Invalid transition
	err = w.TransitionTo(workflow.StateCompleted)
	if err == nil {
		t.Error("TransitionTo(COMPLETED) from PENDING should fail")
	}
}

func TestEntity_SetResult(t *testing.T) {
	w, _ := workflow.New("user-123", "style-456", workflow.Config{
		ImageCount: 2,
		Width:      512,
		Height:     512,
	})

	result := &workflow.Result{
		ImageURLs: []string{"url1", "url2"},
		Duration:  5 * time.Second,
		Error:     "",
	}

	originalUpdatedAt := w.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	w.SetResult(result)

	if w.Result == nil {
		t.Error("Result should not be nil")
	}
	if len(w.Result.ImageURLs) != 2 {
		t.Errorf("Result.ImageURLs length = %d, want 2", len(w.Result.ImageURLs))
	}
	if w.UpdatedAt == originalUpdatedAt {
		t.Error("UpdatedAt should be updated after setting result")
	}
}

func TestEntity_IsFinal(t *testing.T) {
	tests := []struct {
		name  string
		state workflow.State
		want  bool
	}{
		{name: "DRAFT is not final", state: workflow.StateDraft, want: false},
		{name: "PENDING is not final", state: workflow.StatePending, want: false},
		{name: "RUNNING is not final", state: workflow.StateRunning, want: false},
		{name: "PAUSED is not final", state: workflow.StatePaused, want: false},
		{name: "COMPLETED is final", state: workflow.StateCompleted, want: true},
		{name: "FAILED is final", state: workflow.StateFailed, want: true},
		{name: "CANCELLED is final", state: workflow.StateCancelled, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &workflow.Entity{
				UserID:  "user-123",
				StyleID: "style-456",
				State:   tt.state,
			}
			got := w.IsFinal()
			if got != tt.want {
				t.Errorf("IsFinal() = %v, want %v for state %v", got, tt.want, tt.state)
			}
		})
	}
}

func TestEntity_Validate(t *testing.T) {
	validWorkflow, _ := workflow.New("user-123", "style-456", workflow.Config{
		ImageCount: 1,
		Width:      512,
		Height:     512,
	})

	if err := validWorkflow.Validate(); err != nil {
		t.Errorf("Validate() should pass for valid workflow, got %v", err)
	}

	// Test various invalid cases
	invalidCases := []struct {
		name     string
		workflow *workflow.Entity
		want     error
	}{
		{
			name: "empty user ID",
			workflow: &workflow.Entity{
				UserID:  "",
				StyleID: "style-123",
				State:   workflow.StateDraft,
				Config:  workflow.Config{ImageCount: 1, Width: 512, Height: 512},
			},
			want: workflow.ErrEmptyUserID,
		},
		{
			name: "empty style ID",
			workflow: &workflow.Entity{
				UserID:  "user-123",
				StyleID: "",
				State:   workflow.StateDraft,
				Config:  workflow.Config{ImageCount: 1, Width: 512, Height: 512},
			},
			want: workflow.ErrEmptyStyleID,
		},
		{
			name: "invalid state",
			workflow: &workflow.Entity{
				UserID:  "user-123",
				StyleID: "style-456",
				State:   "INVALID",
				Config:  workflow.Config{ImageCount: 1, Width: 512, Height: 512},
			},
			want: workflow.ErrInvalidState,
		},
	}

	for _, tt := range invalidCases {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.workflow.Validate(); err != tt.want {
				t.Errorf("Validate() error = %v, want %v", err, tt.want)
			}
		})
	}
}
