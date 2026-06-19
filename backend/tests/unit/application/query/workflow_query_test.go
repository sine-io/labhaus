package query_test

import (
	"context"
	"testing"

	"github.com/labhaus/backend/internal/application/dto"
	"github.com/labhaus/backend/internal/application/query"
	"github.com/labhaus/backend/internal/domain/workflow"
)

// MockWorkflowRepository is a mock implementation of workflow.Repository
type MockWorkflowRepository struct {
	workflows map[string]*workflow.Entity
}

func NewMockWorkflowRepository() *MockWorkflowRepository {
	repo := &MockWorkflowRepository{
		workflows: make(map[string]*workflow.Entity),
	}
	
	// Prepopulate with test data
	w1, _ := workflow.New("user-1", "style-1", workflow.Config{ImageCount: 1, Width: 512, Height: 512})
	w1.ID = "1"
	w1.State = workflow.StateCompleted
	
	w2, _ := workflow.New("user-1", "style-2", workflow.Config{ImageCount: 2, Width: 1024, Height: 1024})
	w2.ID = "2"
	w2.State = workflow.StateRunning
	
	w3, _ := workflow.New("user-2", "style-1", workflow.Config{ImageCount: 1, Width: 512, Height: 512})
	w3.ID = "3"
	
	repo.workflows["1"] = w1
	repo.workflows["2"] = w2
	repo.workflows["3"] = w3
	
	return repo
}

func (m *MockWorkflowRepository) Create(ctx context.Context, w *workflow.Entity) error {
	return nil
}

func (m *MockWorkflowRepository) FindByID(ctx context.Context, id string) (*workflow.Entity, error) {
	w, ok := m.workflows[id]
	if !ok {
		return nil, workflow.ErrEmptyUserID // Reuse domain error
	}
	return w, nil
}

func (m *MockWorkflowRepository) FindByUserID(ctx context.Context, userID string, filter workflow.Filter) ([]*workflow.Entity, error) {
	result := make([]*workflow.Entity, 0)
	for _, w := range m.workflows {
		if w.UserID == userID {
			// Apply state filter if provided
			if filter.State != nil && w.State != *filter.State {
				continue
			}
			result = append(result, w)
		}
	}
	return result, nil
}

func (m *MockWorkflowRepository) Update(ctx context.Context, w *workflow.Entity) error {
	return nil
}

func (m *MockWorkflowRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *MockWorkflowRepository) UpdateState(ctx context.Context, id string, state workflow.State) error {
	return nil
}

func TestWorkflowQueryHandler_GetWorkflowByID(t *testing.T) {
	repo := NewMockWorkflowRepository()
	handler := query.NewWorkflowQueryHandler(repo)

	result, err := handler.GetWorkflowByID(context.Background(), "1")
	if err != nil {
		t.Fatalf("GetWorkflowByID() error = %v", err)
	}

	if result.ID != "1" {
		t.Errorf("ID = %v, want 1", result.ID)
	}
	if result.UserID != "user-1" {
		t.Errorf("UserID = %v, want user-1", result.UserID)
	}
	if result.State != "COMPLETED" {
		t.Errorf("State = %v, want COMPLETED", result.State)
	}
}

func TestWorkflowQueryHandler_GetWorkflowByID_NotFound(t *testing.T) {
	repo := NewMockWorkflowRepository()
	handler := query.NewWorkflowQueryHandler(repo)

	_, err := handler.GetWorkflowByID(context.Background(), "999")
	if err == nil {
		t.Error("GetWorkflowByID() should return error for non-existent ID")
	}
}

func TestWorkflowQueryHandler_ListUserWorkflows(t *testing.T) {
	repo := NewMockWorkflowRepository()
	handler := query.NewWorkflowQueryHandler(repo)

	filter := dto.WorkflowFilterDTO{
		Limit:  10,
		Offset: 0,
	}

	result, err := handler.ListUserWorkflows(context.Background(), "user-1", filter)
	if err != nil {
		t.Fatalf("ListUserWorkflows() error = %v", err)
	}

	if len(result.Workflows) != 2 {
		t.Errorf("Workflows count = %d, want 2", len(result.Workflows))
	}
	if result.Limit != 10 {
		t.Errorf("Limit = %d, want 10", result.Limit)
	}
}

func TestWorkflowQueryHandler_ListUserWorkflows_WithStateFilter(t *testing.T) {
	repo := NewMockWorkflowRepository()
	handler := query.NewWorkflowQueryHandler(repo)

	runningState := "RUNNING"
	filter := dto.WorkflowFilterDTO{
		State:  &runningState,
		Limit:  10,
		Offset: 0,
	}

	result, err := handler.ListUserWorkflows(context.Background(), "user-1", filter)
	if err != nil {
		t.Fatalf("ListUserWorkflows() error = %v", err)
	}

	if len(result.Workflows) != 1 {
		t.Errorf("Workflows count = %d, want 1", len(result.Workflows))
	}
	if result.Workflows[0].State != "RUNNING" {
		t.Errorf("State = %v, want RUNNING", result.Workflows[0].State)
	}
}

func TestWorkflowQueryHandler_ListUserWorkflows_DefaultLimit(t *testing.T) {
	repo := NewMockWorkflowRepository()
	handler := query.NewWorkflowQueryHandler(repo)

	filter := dto.WorkflowFilterDTO{
		Limit: 0, // Should use default
	}

	result, err := handler.ListUserWorkflows(context.Background(), "user-1", filter)
	if err != nil {
		t.Fatalf("ListUserWorkflows() error = %v", err)
	}

	if result.Limit != 20 {
		t.Errorf("Limit = %d, want 20 (default)", result.Limit)
	}
}

func TestWorkflowQueryHandler_GetWorkflowStatus(t *testing.T) {
	repo := NewMockWorkflowRepository()
	handler := query.NewWorkflowQueryHandler(repo)

	status, err := handler.GetWorkflowStatus(context.Background(), "1")
	if err != nil {
		t.Fatalf("GetWorkflowStatus() error = %v", err)
	}

	if status != "COMPLETED" {
		t.Errorf("Status = %v, want COMPLETED", status)
	}
}

func TestWorkflowQueryHandler_GetWorkflowStatus_NotFound(t *testing.T) {
	repo := NewMockWorkflowRepository()
	handler := query.NewWorkflowQueryHandler(repo)

	_, err := handler.GetWorkflowStatus(context.Background(), "999")
	if err == nil {
		t.Error("GetWorkflowStatus() should return error for non-existent ID")
	}
}
