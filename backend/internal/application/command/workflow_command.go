package command

import (
	"context"
	"errors"
	"time"

	"github.com/labhaus/backend/internal/application/dto"
	"github.com/labhaus/backend/internal/domain/workflow"
)

// WorkflowCommandHandler handles workflow-related commands (write operations)
type WorkflowCommandHandler struct {
	repo workflow.Repository
}

// NewWorkflowCommandHandler creates a new workflow command handler
func NewWorkflowCommandHandler(repo workflow.Repository) *WorkflowCommandHandler {
	return &WorkflowCommandHandler{
		repo: repo,
	}
}

// CreateWorkflow creates a new workflow
func (h *WorkflowCommandHandler) CreateWorkflow(ctx context.Context, userID string, req dto.CreateWorkflowRequest) (*dto.WorkflowDTO, error) {
	// Convert DTO config to domain config
	config := workflow.Config{
		ImageCount: req.Config.ImageCount,
		Width:      req.Config.Width,
		Height:     req.Config.Height,
		Steps:      req.Config.Steps,
		Seed:       req.Config.Seed,
	}

	// Create domain entity
	entity, err := workflow.New(userID, req.StyleID, config)
	if err != nil {
		return nil, err
	}

	// Persist to repository
	if err := h.repo.Create(ctx, entity); err != nil {
		return nil, err
	}

	return toWorkflowDTO(entity), nil
}

// StartWorkflow transitions workflow to PENDING state
func (h *WorkflowCommandHandler) StartWorkflow(ctx context.Context, workflowID string) (*dto.WorkflowDTO, error) {
	entity, err := h.repo.FindByID(ctx, workflowID)
	if err != nil {
		return nil, err
	}

	if err := entity.TransitionTo(workflow.StatePending); err != nil {
		return nil, err
	}

	if err := h.repo.UpdateState(ctx, workflowID, workflow.StatePending); err != nil {
		return nil, err
	}

	return toWorkflowDTO(entity), nil
}

// PauseWorkflow transitions workflow to PAUSED state
func (h *WorkflowCommandHandler) PauseWorkflow(ctx context.Context, workflowID string) (*dto.WorkflowDTO, error) {
	entity, err := h.repo.FindByID(ctx, workflowID)
	if err != nil {
		return nil, err
	}

	if err := entity.TransitionTo(workflow.StatePaused); err != nil {
		return nil, err
	}

	if err := h.repo.UpdateState(ctx, workflowID, workflow.StatePaused); err != nil {
		return nil, err
	}

	return toWorkflowDTO(entity), nil
}

// ResumeWorkflow transitions workflow from PAUSED to RUNNING state
func (h *WorkflowCommandHandler) ResumeWorkflow(ctx context.Context, workflowID string) (*dto.WorkflowDTO, error) {
	entity, err := h.repo.FindByID(ctx, workflowID)
	if err != nil {
		return nil, err
	}

	if err := entity.TransitionTo(workflow.StateRunning); err != nil {
		return nil, err
	}

	if err := h.repo.UpdateState(ctx, workflowID, workflow.StateRunning); err != nil {
		return nil, err
	}

	return toWorkflowDTO(entity), nil
}

// CancelWorkflow transitions workflow to CANCELLED state
func (h *WorkflowCommandHandler) CancelWorkflow(ctx context.Context, workflowID string) (*dto.WorkflowDTO, error) {
	entity, err := h.repo.FindByID(ctx, workflowID)
	if err != nil {
		return nil, err
	}

	if err := entity.TransitionTo(workflow.StateCancelled); err != nil {
		return nil, err
	}

	if err := h.repo.UpdateState(ctx, workflowID, workflow.StateCancelled); err != nil {
		return nil, err
	}

	return toWorkflowDTO(entity), nil
}

// CompleteWorkflow marks workflow as completed with result
func (h *WorkflowCommandHandler) CompleteWorkflow(ctx context.Context, workflowID string, imageURLs []string, duration time.Duration) (*dto.WorkflowDTO, error) {
	entity, err := h.repo.FindByID(ctx, workflowID)
	if err != nil {
		return nil, err
	}

	// Set result
	result := &workflow.Result{
		ImageURLs: imageURLs,
		Duration:  duration,
		Error:     "",
	}
	entity.SetResult(result)

	// Transition to completed
	if err := entity.TransitionTo(workflow.StateCompleted); err != nil {
		return nil, err
	}

	// Update in repository
	if err := h.repo.Update(ctx, entity); err != nil {
		return nil, err
	}

	return toWorkflowDTO(entity), nil
}

// FailWorkflow marks workflow as failed with error message
func (h *WorkflowCommandHandler) FailWorkflow(ctx context.Context, workflowID string, errorMsg string) (*dto.WorkflowDTO, error) {
	entity, err := h.repo.FindByID(ctx, workflowID)
	if err != nil {
		return nil, err
	}

	// Set error result
	result := &workflow.Result{
		ImageURLs: []string{},
		Duration:  0,
		Error:     errorMsg,
	}
	entity.SetResult(result)

	// Transition to failed
	if err := entity.TransitionTo(workflow.StateFailed); err != nil {
		return nil, err
	}

	// Update in repository
	if err := h.repo.Update(ctx, entity); err != nil {
		return nil, err
	}

	return toWorkflowDTO(entity), nil
}

// DeleteWorkflow deletes a workflow by ID
func (h *WorkflowCommandHandler) DeleteWorkflow(ctx context.Context, workflowID string) error {
	// Check if workflow exists
	_, err := h.repo.FindByID(ctx, workflowID)
	if err != nil {
		return err
	}

	return h.repo.Delete(ctx, workflowID)
}

// toWorkflowDTO converts domain entity to DTO
func toWorkflowDTO(entity *workflow.Entity) *dto.WorkflowDTO {
	dtoWorkflow := &dto.WorkflowDTO{
		ID:      entity.ID,
		UserID:  entity.UserID,
		StyleID: entity.StyleID,
		State:   string(entity.State),
		Config: dto.WorkflowConfig{
			ImageCount: entity.Config.ImageCount,
			Width:      entity.Config.Width,
			Height:     entity.Config.Height,
			Steps:      entity.Config.Steps,
			Seed:       entity.Config.Seed,
		},
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	// Convert result if present
	if entity.Result != nil {
		dtoWorkflow.Result = &dto.WorkflowResult{
			ImageURLs: entity.Result.ImageURLs,
			Duration:  entity.Result.Duration.Milliseconds(),
			Error:     entity.Result.Error,
		}
	}

	return dtoWorkflow
}

// Common errors
var (
	ErrWorkflowNotFound = errors.New("workflow not found")
)
