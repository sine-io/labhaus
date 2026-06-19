package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/labhaus/backend/internal/domain/workflow"
	"gorm.io/gorm"
)

// WorkflowRepository implements workflow.Repository using GORM
type WorkflowRepository struct {
	db *gorm.DB
}

// NewWorkflowRepository creates a new workflow repository
func NewWorkflowRepository(db *gorm.DB) *WorkflowRepository {
	return &WorkflowRepository{db: db}
}

// Create saves a new workflow
func (r *WorkflowRepository) Create(ctx context.Context, w *workflow.Entity) error {
	// Generate ID if not set
	if w.ID == "" {
		w.ID = uuid.New().String()
	}

	model, err := r.toModel(w)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Create(model).Error
}

// FindByID retrieves a workflow by ID
func (r *WorkflowRepository) FindByID(ctx context.Context, id string) (*workflow.Entity, error) {
	var model WorkflowModel
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("workflow not found")
		}
		return nil, err
	}

	return r.toEntity(&model)
}

// FindByUserID retrieves all workflows for a user
func (r *WorkflowRepository) FindByUserID(ctx context.Context, userID string, filter workflow.Filter) ([]*workflow.Entity, error) {
	var models []WorkflowModel
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	// Apply state filter if provided
	if filter.State != nil {
		query = query.Where("state = ?", string(*filter.State))
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
	entities := make([]*workflow.Entity, len(models))
	for i, model := range models {
		entity, err := r.toEntity(&model)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}

	return entities, nil
}

// Update modifies an existing workflow
func (r *WorkflowRepository) Update(ctx context.Context, w *workflow.Entity) error {
	model, err := r.toModel(w)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Model(&WorkflowModel{}).Where("id = ?", w.ID).Updates(model)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("workflow not found")
	}

	return nil
}

// Delete removes a workflow by ID (soft delete)
func (r *WorkflowRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&WorkflowModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("workflow not found")
	}

	return nil
}

// UpdateState updates only the workflow state
func (r *WorkflowRepository) UpdateState(ctx context.Context, id string, state workflow.State) error {
	result := r.db.WithContext(ctx).Model(&WorkflowModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"state":      string(state),
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("workflow not found")
	}

	return nil
}

// toModel converts domain entity to database model
func (r *WorkflowRepository) toModel(w *workflow.Entity) (*WorkflowModel, error) {
	// Convert config to JSONB
	configJSON, err := json.Marshal(w.Config)
	if err != nil {
		return nil, err
	}
	var configMap JSONB
	if err := json.Unmarshal(configJSON, &configMap); err != nil {
		return nil, err
	}

	model := &WorkflowModel{
		ID:        w.ID,
		UserID:    w.UserID,
		StyleID:   w.StyleID,
		State:     string(w.State),
		Config:    configMap,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}

	// Convert result to JSONB if present
	if w.Result != nil {
		resultJSON, err := json.Marshal(w.Result)
		if err != nil {
			return nil, err
		}
		var resultMap JSONB
		if err := json.Unmarshal(resultJSON, &resultMap); err != nil {
			return nil, err
		}
		model.Result = resultMap
	}

	return model, nil
}

// toEntity converts database model to domain entity
func (r *WorkflowRepository) toEntity(m *WorkflowModel) (*workflow.Entity, error) {
	// Convert config from JSONB
	configJSON, err := json.Marshal(m.Config)
	if err != nil {
		return nil, err
	}
	var config workflow.Config
	if err := json.Unmarshal(configJSON, &config); err != nil {
		return nil, err
	}

	entity := &workflow.Entity{
		ID:        m.ID,
		UserID:    m.UserID,
		StyleID:   m.StyleID,
		State:     workflow.State(m.State),
		Config:    config,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	// Convert result from JSONB if present
	if m.Result != nil {
		resultJSON, err := json.Marshal(m.Result)
		if err != nil {
			return nil, err
		}
		var result workflow.Result
		if err := json.Unmarshal(resultJSON, &result); err != nil {
			return nil, err
		}
		entity.Result = &result
	}

	return entity, nil
}
