package query

import (
	"context"

	"github.com/labhaus/backend/internal/application/dto"
	"github.com/labhaus/backend/internal/domain/style"
)

// StyleQueryHandler handles style-related queries (read operations)
type StyleQueryHandler struct {
	repo style.Repository
}

// NewStyleQueryHandler creates a new style query handler
func NewStyleQueryHandler(repo style.Repository) *StyleQueryHandler {
	return &StyleQueryHandler{
		repo: repo,
	}
}

// GetStyleByID retrieves a style by ID
func (h *StyleQueryHandler) GetStyleByID(ctx context.Context, id string) (*dto.StyleDTO, error) {
	entity, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toStyleDTO(entity), nil
}

// ListStyles retrieves all styles with optional filters
func (h *StyleQueryHandler) ListStyles(ctx context.Context, filter dto.StyleFilterDTO) (*dto.StyleListResponse, error) {
	// Convert DTO filter to domain filter
	domainFilter := style.Filter{
		Category: filter.Category,
		Tags:     filter.Tags,
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	}

	// Set default limit if not provided
	if domainFilter.Limit <= 0 {
		domainFilter.Limit = 20
	}

	entities, err := h.repo.FindAll(ctx, domainFilter)
	if err != nil {
		return nil, err
	}

	// Convert entities to DTOs
	styleDTOs := make([]dto.StyleDTO, len(entities))
	for i, entity := range entities {
		styleDTOs[i] = *toStyleDTO(entity)
	}

	return &dto.StyleListResponse{
		Styles: styleDTOs,
		Total:  len(styleDTOs),
		Limit:  domainFilter.Limit,
		Offset: domainFilter.Offset,
	}, nil
}

// SearchStyles performs full-text search on styles
func (h *StyleQueryHandler) SearchStyles(ctx context.Context, query string, limit int) ([]dto.StyleDTO, error) {
	if limit <= 0 {
		limit = 20
	}

	entities, err := h.repo.Search(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	// Convert entities to DTOs
	styleDTOs := make([]dto.StyleDTO, len(entities))
	for i, entity := range entities {
		styleDTOs[i] = *toStyleDTO(entity)
	}

	return styleDTOs, nil
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
