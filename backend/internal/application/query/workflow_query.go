package query

import (
	"context"

	"github.com/labhaus/backend/internal/application/dto"
	"github.com/labhaus/backend/internal/domain/workflow"
)

// WorkflowQueryHandler handles workflow-related queries (read operations)
type WorkflowQueryHandler struct {
	repo workflow.Repository
}

// NewWorkflowQueryHandler creates a new workflow query handler
func NewWorkflowQueryHandler(repo workflow.Repository) *WorkflowQueryHandler {
	return &WorkflowQueryHandler{
		repo: repo,
	}
}

// GetWorkflowByID retrieves a workflow by ID
func (h *WorkflowQueryHandler) GetWorkflowByID(ctx context.Context, id string) (*dto.WorkflowDTO, error) {
	entity, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toWorkflowDTO(entity), nil
}

// ListUserWorkflows retrieves all workflows for a user with optional filters
func (h *WorkflowQueryHandler) ListUserWorkflows(ctx context.Context, userID string, filter dto.WorkflowFilterDTO) (*dto.WorkflowListResponse, error) {
	// Convert DTO filter to domain filter
	domainFilter := workflow.Filter{
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}

	// Set state filter if provided
	if filter.State != nil && *filter.State != "" {
		state := workflow.State(*filter.State)
		domainFilter.State = &state
	}

	// Set default limit if not provided
	if domainFilter.Limit <= 0 {
		domainFilter.Limit = 20
	}

	entities, err := h.repo.FindByUserID(ctx, userID, domainFilter)
	if err != nil {
		return nil, err
	}

	// Convert entities to DTOs
	workflowDTOs := make([]dto.WorkflowDTO, len(entities))
	for i, entity := range entities {
		workflowDTOs[i] = *toWorkflowDTO(entity)
	}

	return &dto.WorkflowListResponse{
		Workflows: workflowDTOs,
		Total:     len(workflowDTOs),
		Limit:     domainFilter.Limit,
		Offset:    domainFilter.Offset,
	}, nil
}

// GetWorkflowStatus retrieves only the status of a workflow
func (h *WorkflowQueryHandler) GetWorkflowStatus(ctx context.Context, id string) (string, error) {
	entity, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return "", err
	}

	return string(entity.State), nil
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
