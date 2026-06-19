package command

import (
	"context"
	"errors"

	"github.com/labhaus/backend/internal/application/dto"
	"github.com/labhaus/backend/internal/domain/user"
)

// UserCommandHandler handles user-related commands (write operations)
type UserCommandHandler struct {
	repo           user.Repository
	passwordHasher PasswordHasher
}

// PasswordHasher interface for password hashing (DIP)
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

// NewUserCommandHandler creates a new user command handler
func NewUserCommandHandler(repo user.Repository, hasher PasswordHasher) *UserCommandHandler {
	return &UserCommandHandler{
		repo:           repo,
		passwordHasher: hasher,
	}
}

// RegisterUser registers a new user
func (h *UserCommandHandler) RegisterUser(ctx context.Context, req dto.RegisterUserRequest) (*dto.UserDTO, error) {
	// Check if email already exists
	exists, err := h.repo.Exists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := h.passwordHasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	// Create domain entity
	entity, err := user.New(req.Email, hashedPassword, req.Name, user.RoleUser)
	if err != nil {
		return nil, err
	}

	// Persist to repository
	if err := h.repo.Create(ctx, entity); err != nil {
		return nil, err
	}

	return toUserDTO(entity), nil
}

// UpdateUser updates user information
func (h *UserCommandHandler) UpdateUser(ctx context.Context, userID string, req dto.UpdateUserRequest) (*dto.UserDTO, error) {
	// Find existing entity
	entity, err := h.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update entity
	if err := entity.Update(req.Email, req.Name, entity.Role); err != nil {
		return nil, err
	}

	// Persist changes
	if err := h.repo.Update(ctx, entity); err != nil {
		return nil, err
	}

	return toUserDTO(entity), nil
}

// ChangePassword changes user password
func (h *UserCommandHandler) ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) error {
	// Find existing entity
	entity, err := h.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := h.passwordHasher.Compare(entity.PasswordHash, req.OldPassword); err != nil {
		return ErrInvalidPassword
	}

	// Hash new password
	hashedPassword, err := h.passwordHasher.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password hash
	entity.PasswordHash = hashedPassword

	// Persist changes
	return h.repo.Update(ctx, entity)
}

// DeleteUser deletes a user by ID
func (h *UserCommandHandler) DeleteUser(ctx context.Context, userID string) error {
	// Check if user exists
	_, err := h.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	return h.repo.Delete(ctx, userID)
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

// Common errors
var (
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidPassword     = errors.New("invalid password")
)
