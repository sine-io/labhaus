package command

import (
	"context"
	"errors"

	"github.com/labhaus/backend/internal/application/dto"
	"github.com/labhaus/backend/internal/domain/style"
)

// StyleCommandHandler handles style-related commands (write operations)
type StyleCommandHandler struct {
	repo style.Repository
}

// NewStyleCommandHandler creates a new style command handler
func NewStyleCommandHandler(repo style.Repository) *StyleCommandHandler {
	return &StyleCommandHandler{
		repo: repo,
	}
}

// CreateStyle creates a new style
func (h *StyleCommandHandler) CreateStyle(ctx context.Context, req dto.CreateStyleRequest) (*dto.StyleDTO, error) {
	// Create domain entity
	entity, err := style.New(req.Name, req.Description, req.Prompt, req.Category, req.Tags)
	if err != nil {
		return nil, err
	}

	// Persist to repository
	if err := h.repo.Create(ctx, entity); err != nil {
		return nil, err
	}

	// Convert to DTO
	return toStyleDTO(entity), nil
}

// UpdateStyle updates an existing style
func (h *StyleCommandHandler) UpdateStyle(ctx context.Context, id string, req dto.UpdateStyleRequest) (*dto.StyleDTO, error) {
	// Find existing entity
	entity, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update entity
	if err := entity.Update(req.Name, req.Description, req.Prompt, req.Category, req.Tags); err != nil {
		return nil, err
	}

	// Persist changes
	if err := h.repo.Update(ctx, entity); err != nil {
		return nil, err
	}

	return toStyleDTO(entity), nil
}

// DeleteStyle deletes a style by ID
func (h *StyleCommandHandler) DeleteStyle(ctx context.Context, id string) error {
	// Check if style exists
	_, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete from repository
	return h.repo.Delete(ctx, id)
}

// toStyleDTO converts domain entity to DTO
func toStyleDTO(entity *style.Entity) *dto.StyleDTO {
	return &dto.StyleDTO{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Prompt:      entity.Prompt,
		Category:    entity.Category,
		Tags:        entity.Tags,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

// Common errors
var (
	ErrStyleNotFound = errors.New("style not found")
)
