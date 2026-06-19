package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/labhaus/backend/internal/domain/user"
	"gorm.io/gorm"
)

// UserRepository implements user.Repository using GORM
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create saves a new user
func (r *UserRepository) Create(ctx context.Context, u *user.Entity) error {
	// Generate ID if not set
	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	model := r.toModel(u)
	return r.db.WithContext(ctx).Create(model).Error
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.Entity, error) {
	var model UserModel
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return r.toEntity(&model), nil
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.Entity, error) {
	var model UserModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return r.toEntity(&model), nil
}

// Update modifies an existing user
func (r *UserRepository) Update(ctx context.Context, u *user.Entity) error {
	model := r.toModel(u)
	result := r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", u.ID).Updates(model)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Delete removes a user by ID (soft delete)
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&UserModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Exists checks if a user with the given email exists
func (r *UserRepository) Exists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&UserModel{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// toModel converts domain entity to database model
func (r *UserRepository) toModel(u *user.Entity) *UserModel {
	return &UserModel{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Name:         u.Name,
		Role:         string(u.Role),
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// toEntity converts database model to domain entity
func (r *UserRepository) toEntity(m *UserModel) *user.Entity {
	return &user.Entity{
		ID:           m.ID,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		Name:         m.Name,
		Role:         user.Role(m.Role),
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
