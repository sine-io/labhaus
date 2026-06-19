package query

import (
	"context"

	"github.com/labhaus/backend/internal/application/dto"
	"github.com/labhaus/backend/internal/domain/user"
)

// UserQueryHandler handles user-related queries (read operations)
type UserQueryHandler struct {
	repo user.Repository
}

// NewUserQueryHandler creates a new user query handler
func NewUserQueryHandler(repo user.Repository) *UserQueryHandler {
	return &UserQueryHandler{
		repo: repo,
	}
}

// GetUserByID retrieves a user by ID
func (h *UserQueryHandler) GetUserByID(ctx context.Context, id string) (*dto.UserDTO, error) {
	entity, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toUserDTO(entity), nil
}

// GetUserByEmail retrieves a user by email
func (h *UserQueryHandler) GetUserByEmail(ctx context.Context, email string) (*dto.UserDTO, error) {
	entity, err := h.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return toUserDTO(entity), nil
}

// CheckEmailExists checks if an email is already registered
func (h *UserQueryHandler) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	return h.repo.Exists(ctx, email)
}

// toUserDTO converts domain entity to DTO
func toUserDTO(entity *user.Entity) *dto.UserDTO {
	return &dto.UserDTO{
		ID:        entity.ID,
		Email:     entity.Email,
		Name:      entity.Name,
		Role:      string(entity.Role),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
