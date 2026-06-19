package persistence

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/labhaus/backend/internal/domain/style"
	"gorm.io/gorm"
)

// StyleRepository implements style.Repository using GORM
type StyleRepository struct {
	db *gorm.DB
}

// NewStyleRepository creates a new style repository
func NewStyleRepository(db *gorm.DB) *StyleRepository {
	return &StyleRepository{db: db}
}

// Create saves a new style
func (r *StyleRepository) Create(ctx context.Context, s *style.Entity) error {
	// Generate ID if not set
	if s.ID == "" {
		s.ID = uuid.New().String()
	}

	model, err := r.toModel(s)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Create(model).Error
}

// FindByID retrieves a style by ID
func (r *StyleRepository) FindByID(ctx context.Context, id string) (*style.Entity, error) {
	var model StyleModel
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("style not found")
		}
		return nil, err
	}

	return r.toEntity(&model)
}

// FindAll retrieves all styles with optional filters
func (r *StyleRepository) FindAll(ctx context.Context, filter style.Filter) ([]*style.Entity, error) {
	var models []StyleModel
	query := r.db.WithContext(ctx).Model(&StyleModel{})

	// Apply filters
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	// Tags filter (simple contains check)
	if len(filter.Tags) > 0 {
		for _, tag := range filter.Tags {
			query = query.Where("tags LIKE ?", "%"+tag+"%")
		}
	}

	// Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	err := query.Order("created_at DESC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	// Convert to entities
	entities := make([]*style.Entity, len(models))
	for i, model := range models {
		entity, err := r.toEntity(&model)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}

	return entities, nil
}

// Update modifies an existing style
func (r *StyleRepository) Update(ctx context.Context, s *style.Entity) error {
	model, err := r.toModel(s)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Model(&StyleModel{}).Where("id = ?", s.ID).Updates(model)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("style not found")
	}

	return nil
}

// Delete removes a style by ID (soft delete)
func (r *StyleRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&StyleModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("style not found")
	}

	return nil
}

// Search performs full-text search on styles
func (r *StyleRepository) Search(ctx context.Context, query string, limit int) ([]*style.Entity, error) {
	var models []StyleModel

	// Simple LIKE search on name, description, and prompt
	searchPattern := "%" + query + "%"
	err := r.db.WithContext(ctx).
		Where("name LIKE ? OR description LIKE ? OR prompt LIKE ?", searchPattern, searchPattern, searchPattern).
		Limit(limit).
		Order("created_at DESC").
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	// Convert to entities
	entities := make([]*style.Entity, len(models))
	for i, model := range models {
		entity, err := r.toEntity(&model)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}

	return entities, nil
}

// toModel converts domain entity to database model
func (r *StyleRepository) toModel(s *style.Entity) (*StyleModel, error) {
	// Marshal tags to JSON
	tagsJSON, err := json.Marshal(s.Tags)
	if err != nil {
		return nil, err
	}

	return &StyleModel{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Prompt:      s.Prompt,
		Category:    s.Category,
		Tags:        string(tagsJSON),
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}, nil
}

// toEntity converts database model to domain entity
func (r *StyleRepository) toEntity(m *StyleModel) (*style.Entity, error) {
	// Unmarshal tags from JSON
	var tags []string
	if m.Tags != "" {
		if err := json.Unmarshal([]byte(m.Tags), &tags); err != nil {
			return nil, err
		}
	}

	return &style.Entity{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Prompt:      m.Prompt,
		Category:    m.Category,
		Tags:        tags,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}, nil
}
