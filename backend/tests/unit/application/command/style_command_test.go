package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/labhaus/backend/internal/application/command"
	"github.com/labhaus/backend/internal/application/dto"
	"github.com/labhaus/backend/internal/domain/style"
)

// MockStyleRepository is a mock implementation of style.Repository
type MockStyleRepository struct {
	styles map[string]*style.Entity
	nextID int
}

func NewMockStyleRepository() *MockStyleRepository {
	return &MockStyleRepository{
		styles: make(map[string]*style.Entity),
		nextID: 1,
	}
}

func (m *MockStyleRepository) Create(ctx context.Context, s *style.Entity) error {
	id := string(rune(m.nextID + 48)) // Simple ID generation
	m.nextID++
	s.ID = id
	m.styles[id] = s
	return nil
}

func (m *MockStyleRepository) FindByID(ctx context.Context, id string) (*style.Entity, error) {
	s, ok := m.styles[id]
	if !ok {
		return nil, errors.New("style not found")
	}
	return s, nil
}

func (m *MockStyleRepository) FindAll(ctx context.Context, filter style.Filter) ([]*style.Entity, error) {
	result := make([]*style.Entity, 0, len(m.styles))
	for _, s := range m.styles {
		result = append(result, s)
	}
	return result, nil
}

func (m *MockStyleRepository) Update(ctx context.Context, s *style.Entity) error {
	if _, ok := m.styles[s.ID]; !ok {
		return errors.New("style not found")
	}
	m.styles[s.ID] = s
	return nil
}

func (m *MockStyleRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.styles[id]; !ok {
		return errors.New("style not found")
	}
	delete(m.styles, id)
	return nil
}

func (m *MockStyleRepository) Search(ctx context.Context, query string, limit int) ([]*style.Entity, error) {
	result := make([]*style.Entity, 0)
	for _, s := range m.styles {
		result = append(result, s)
		if len(result) >= limit {
			break
		}
	}
	return result, nil
}

func TestStyleCommandHandler_CreateStyle(t *testing.T) {
	repo := NewMockStyleRepository()
	handler := command.NewStyleCommandHandler(repo)

	req := dto.CreateStyleRequest{
		Name:        "Anime Style",
		Description: "Japanese anime art",
		Prompt:      "anime, colorful, expressive",
		Category:    "Art",
		Tags:        []string{"anime", "japanese"},
	}

	result, err := handler.CreateStyle(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateStyle() error = %v", err)
	}

	if result.Name != req.Name {
		t.Errorf("Name = %v, want %v", result.Name, req.Name)
	}
	if result.ID == "" {
		t.Error("ID should not be empty")
	}
}

func TestStyleCommandHandler_UpdateStyle(t *testing.T) {
	repo := NewMockStyleRepository()
	handler := command.NewStyleCommandHandler(repo)

	// Create initial style
	createReq := dto.CreateStyleRequest{
		Name:        "Original",
		Description: "Original desc",
		Prompt:      "original prompt",
		Category:    "Cat",
		Tags:        []string{"tag1"},
	}
	created, _ := handler.CreateStyle(context.Background(), createReq)

	// Update style
	updateReq := dto.UpdateStyleRequest{
		Name:        "Updated",
		Description: "Updated desc",
		Prompt:      "updated prompt",
		Category:    "NewCat",
		Tags:        []string{"tag2"},
	}

	result, err := handler.UpdateStyle(context.Background(), created.ID, updateReq)
	if err != nil {
		t.Fatalf("UpdateStyle() error = %v", err)
	}

	if result.Name != "Updated" {
		t.Errorf("Name = %v, want Updated", result.Name)
	}
	if result.Description != "Updated desc" {
		t.Errorf("Description = %v, want Updated desc", result.Description)
	}
}

func TestStyleCommandHandler_DeleteStyle(t *testing.T) {
	repo := NewMockStyleRepository()
	handler := command.NewStyleCommandHandler(repo)

	// Create style
	createReq := dto.CreateStyleRequest{
		Name:   "Test",
		Prompt: "test prompt",
	}
	created, _ := handler.CreateStyle(context.Background(), createReq)

	// Delete style
	err := handler.DeleteStyle(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("DeleteStyle() error = %v", err)
	}

	// Verify deletion
	_, err = repo.FindByID(context.Background(), created.ID)
	if err == nil {
		t.Error("Style should be deleted")
	}
}

func TestStyleCommandHandler_CreateStyle_ValidationError(t *testing.T) {
	repo := NewMockStyleRepository()
	handler := command.NewStyleCommandHandler(repo)

	// Empty name should fail
	req := dto.CreateStyleRequest{
		Name:   "",
		Prompt: "test prompt",
	}

	_, err := handler.CreateStyle(context.Background(), req)
	if err == nil {
		t.Error("CreateStyle() should fail with empty name")
	}
}
