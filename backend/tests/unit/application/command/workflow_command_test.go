package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/labhaus/backend/internal/application/command"
	"github.com/labhaus/backend/internal/application/dto"
	"github.com/labhaus/backend/internal/domain/workflow"
)

// MockWorkflowRepository is a mock implementation of workflow.Repository
type MockWorkflowRepository struct {
	workflows map[string]*workflow.Entity
	nextID    int
}

func NewMockWorkflowRepository() *MockWorkflowRepository {
	return &MockWorkflowRepository{
		workflows: make(map[string]*workflow.Entity),
		nextID:    1,
	}
}

func (m *MockWorkflowRepository) Create(ctx context.Context, w *workflow.Entity) error {
	id := string(rune(m.nextID + 48))
	m.nextID++
	w.ID = id
	m.workflows[id] = w
	return nil
}

func (m *MockWorkflowRepository) FindByID(ctx context.Context, id string) (*workflow.Entity, error) {
	w, ok := m.workflows[id]
	if !ok {
		return nil, errors.New("workflow not found")
	}
	return w, nil
}

func (m *MockWorkflowRepository) FindByUserID(ctx context.Context, userID string, filter workflow.Filter) ([]*workflow.Entity, error) {
	result := make([]*workflow.Entity, 0)
	for _, w := range m.workflows {
		if w.UserID == userID {
			result = append(result, w)
		}
	}
	return result, nil
}

func (m *MockWorkflowRepository) Update(ctx context.Context, w *workflow.Entity) error {
	if _, ok := m.workflows[w.ID]; !ok {
		return errors.New("workflow not found")
	}
	m.workflows[w.ID] = w
	return nil
}

func (m *MockWorkflowRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.workflows[id]; !ok {
		return errors.New("workflow not found")
	}
	delete(m.workflows, id)
	return nil
}

func (m *MockWorkflowRepository) UpdateState(ctx context.Context, id string, state workflow.State) error {
	w, ok := m.workflows[id]
	if !ok {
		return errors.New("workflow not found")
	}
	w.State = state
	w.UpdatedAt = time.Now()
	return nil
}

func TestWorkflowCommandHandler_CreateWorkflow(t *testing.T) {
	repo := NewMockWorkflowRepository()
	handler := command.NewWorkflowCommandHandler(repo)

	req := dto.CreateWorkflowRequest{
		StyleID: "style-123",
		Config: dto.WorkflowConfig{
			ImageCount: 4,
			Width:      512,
			Height:     512,
			Steps:      30,
			Seed:       12345,
		},
	}

	result, err := handler.CreateWorkflow(context.Background(), "user-123", req)
	if err != nil {
		t.Fatalf("CreateWorkflow() error = %v", err)
	}

	if result.UserID != "user-123" {
		t.Errorf("UserID = %v, want user-123", result.UserID)
	}
	if result.StyleID != "style-123" {
		t.Errorf("StyleID = %v, want style-123", result.StyleID)
	}
	if result.State != "DRAFT" {
		t.Errorf("Initial state = %v, want DRAFT", result.State)
	}
}

func TestWorkflowCommandHandler_StateTransitions(t *testing.T) {
	repo := NewMockWorkflowRepository()
	handler := command.NewWorkflowCommandHandler(repo)

	// Create workflow
	req := dto.CreateWorkflowRequest{
		StyleID: "style-123",
		Config: dto.WorkflowConfig{
			ImageCount: 1,
			Width:      512,
			Height:     512,
		},
	}
	created, _ := handler.CreateWorkflow(context.Background(), "user-123", req)

	// Test valid transitions
	t.Run("DRAFT -> PENDING", func(t *testing.T) {
		result, err := handler.StartWorkflow(context.Background(), created.ID)
		if err != nil {
			t.Fatalf("StartWorkflow() error = %v", err)
		}
		if result.State != "PENDING" {
			t.Errorf("State = %v, want PENDING", result.State)
		}
	})

	// Set to RUNNING manually for pause test
	repo.UpdateState(context.Background(), created.ID, workflow.StateRunning)

	t.Run("RUNNING -> PAUSED", func(t *testing.T) {
		result, err := handler.PauseWorkflow(context.Background(), created.ID)
		if err != nil {
			t.Fatalf("PauseWorkflow() error = %v", err)
		}
		if result.State != "PAUSED" {
			t.Errorf("State = %v, want PAUSED", result.State)
		}
	})

	t.Run("PAUSED -> RUNNING", func(t *testing.T) {
		result, err := handler.ResumeWorkflow(context.Background(), created.ID)
		if err != nil {
			t.Fatalf("ResumeWorkflow() error = %v", err)
		}
		if result.State != "RUNNING" {
			t.Errorf("State = %v, want RUNNING", result.State)
		}
	})
}

func TestWorkflowCommandHandler_CompleteWorkflow(t *testing.T) {
	repo := NewMockWorkflowRepository()
	handler := command.NewWorkflowCommandHandler(repo)

	// Create and start workflow
	req := dto.CreateWorkflowRequest{
		StyleID: "style-123",
		Config: dto.WorkflowConfig{
			ImageCount: 2,
			Width:      512,
			Height:     512,
		},
	}
	created, _ := handler.CreateWorkflow(context.Background(), "user-123", req)
	
	// Set to RUNNING
	repo.UpdateState(context.Background(), created.ID, workflow.StateRunning)

	// Complete workflow
	imageURLs := []string{"url1", "url2"}
	duration := 5 * time.Second

	result, err := handler.CompleteWorkflow(context.Background(), created.ID, imageURLs, duration)
	if err != nil {
		t.Fatalf("CompleteWorkflow() error = %v", err)
	}

	if result.State != "COMPLETED" {
		t.Errorf("State = %v, want COMPLETED", result.State)
	}
	if result.Result == nil {
		t.Fatal("Result should not be nil")
	}
	if len(result.Result.ImageURLs) != 2 {
		t.Errorf("ImageURLs length = %d, want 2", len(result.Result.ImageURLs))
	}
	if result.Result.Duration != 5000 {
		t.Errorf("Duration = %d ms, want 5000", result.Result.Duration)
	}
}

func TestWorkflowCommandHandler_FailWorkflow(t *testing.T) {
	repo := NewMockWorkflowRepository()
	handler := command.NewWorkflowCommandHandler(repo)

	// Create and start workflow
	req := dto.CreateWorkflowRequest{
		StyleID: "style-123",
		Config: dto.WorkflowConfig{
			ImageCount: 1,
			Width:      512,
			Height:     512,
		},
	}
	created, _ := handler.CreateWorkflow(context.Background(), "user-123", req)
	
	// Set to RUNNING
	repo.UpdateState(context.Background(), created.ID, workflow.StateRunning)

	// Fail workflow
	errorMsg := "Image generation failed"
	result, err := handler.FailWorkflow(context.Background(), created.ID, errorMsg)
	if err != nil {
		t.Fatalf("FailWorkflow() error = %v", err)
	}

	if result.State != "FAILED" {
		t.Errorf("State = %v, want FAILED", result.State)
	}
	if result.Result == nil {
		t.Fatal("Result should not be nil")
	}
	if result.Result.Error != errorMsg {
		t.Errorf("Error = %v, want %v", result.Result.Error, errorMsg)
	}
}

func TestWorkflowCommandHandler_CancelWorkflow(t *testing.T) {
	repo := NewMockWorkflowRepository()
	handler := command.NewWorkflowCommandHandler(repo)

	// Create workflow
	req := dto.CreateWorkflowRequest{
		StyleID: "style-123",
		Config: dto.WorkflowConfig{
			ImageCount: 1,
			Width:      512,
			Height:     512,
		},
	}
	created, _ := handler.CreateWorkflow(context.Background(), "user-123", req)

	// Cancel from DRAFT
	result, err := handler.CancelWorkflow(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("CancelWorkflow() error = %v", err)
	}

	if result.State != "CANCELLED" {
		t.Errorf("State = %v, want CANCELLED", result.State)
	}
}

func TestWorkflowCommandHandler_DeleteWorkflow(t *testing.T) {
	repo := NewMockWorkflowRepository()
	handler := command.NewWorkflowCommandHandler(repo)

	// Create workflow
	req := dto.CreateWorkflowRequest{
		StyleID: "style-123",
		Config: dto.WorkflowConfig{
			ImageCount: 1,
			Width:      512,
			Height:     512,
		},
	}
	created, _ := handler.CreateWorkflow(context.Background(), "user-123", req)

	// Delete workflow
	err := handler.DeleteWorkflow(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("DeleteWorkflow() error = %v", err)
	}

	// Verify deletion
	_, err = repo.FindByID(context.Background(), created.ID)
	if err == nil {
		t.Error("Workflow should be deleted")
	}
}
