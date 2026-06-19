package query_test

import (
	"context"
	"testing"

	"github.com/labhaus/backend/internal/application/dto"
	"github.com/labhaus/backend/internal/application/query"
	"github.com/labhaus/backend/internal/domain/style"
)

// MockStyleRepository is a mock implementation of style.Repository
type MockStyleRepository struct {
	styles map[string]*style.Entity
}

func NewMockStyleRepository() *MockStyleRepository {
	repo := &MockStyleRepository{
		styles: make(map[string]*style.Entity),
	}
	
	// Prepopulate with test data
	s1, _ := style.New("Anime", "Anime style", "anime art", "Art", []string{"anime"})
	s1.ID = "1"
	s2, _ := style.New("Realistic", "Photo realistic", "realistic photo", "Photo", []string{"realistic"})
	s2.ID = "2"
	
	repo.styles["1"] = s1
	repo.styles["2"] = s2
	
	return repo
}

func (m *MockStyleRepository) Create(ctx context.Context, s *style.Entity) error {
	return nil
}

func (m *MockStyleRepository) FindByID(ctx context.Context, id string) (*style.Entity, error) {
	s, ok := m.styles[id]
	if !ok {
		return nil, style.ErrEmptyName // Reuse domain error
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
	return nil
}

func (m *MockStyleRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *MockStyleRepository) Search(ctx context.Context, q string, limit int) ([]*style.Entity, error) {
	result := make([]*style.Entity, 0)
	for _, s := range m.styles {
		result = append(result, s)
		if len(result) >= limit {
			break
		}
	}
	return result, nil
}

func TestStyleQueryHandler_GetStyleByID(t *testing.T) {
	repo := NewMockStyleRepository()
	handler := query.NewStyleQueryHandler(repo)

	result, err := handler.GetStyleByID(context.Background(), "1")
	if err != nil {
		t.Fatalf("GetStyleByID() error = %v", err)
	}

	if result.ID != "1" {
		t.Errorf("ID = %v, want 1", result.ID)
	}
	if result.Name != "Anime" {
		t.Errorf("Name = %v, want Anime", result.Name)
	}
}

func TestStyleQueryHandler_GetStyleByID_NotFound(t *testing.T) {
	repo := NewMockStyleRepository()
	handler := query.NewStyleQueryHandler(repo)

	_, err := handler.GetStyleByID(context.Background(), "999")
	if err == nil {
		t.Error("GetStyleByID() should return error for non-existent ID")
	}
}

func TestStyleQueryHandler_ListStyles(t *testing.T) {
	repo := NewMockStyleRepository()
	handler := query.NewStyleQueryHandler(repo)

	filter := dto.StyleFilterDTO{
		Limit:  10,
		Offset: 0,
	}

	result, err := handler.ListStyles(context.Background(), filter)
	if err != nil {
		t.Fatalf("ListStyles() error = %v", err)
	}

	if len(result.Styles) != 2 {
		t.Errorf("Styles count = %d, want 2", len(result.Styles))
	}
	if result.Limit != 10 {
		t.Errorf("Limit = %d, want 10", result.Limit)
	}
}

func TestStyleQueryHandler_ListStyles_DefaultLimit(t *testing.T) {
	repo := NewMockStyleRepository()
	handler := query.NewStyleQueryHandler(repo)

	filter := dto.StyleFilterDTO{
		Limit: 0, // Should use default
	}

	result, err := handler.ListStyles(context.Background(), filter)
	if err != nil {
		t.Fatalf("ListStyles() error = %v", err)
	}

	if result.Limit != 20 {
		t.Errorf("Limit = %d, want 20 (default)", result.Limit)
	}
}

func TestStyleQueryHandler_SearchStyles(t *testing.T) {
	repo := NewMockStyleRepository()
	handler := query.NewStyleQueryHandler(repo)

	results, err := handler.SearchStyles(context.Background(), "anime", 10)
	if err != nil {
		t.Fatalf("SearchStyles() error = %v", err)
	}

	if len(results) == 0 {
		t.Error("SearchStyles() should return results")
	}
}

func TestStyleQueryHandler_SearchStyles_DefaultLimit(t *testing.T) {
	repo := NewMockStyleRepository()
	handler := query.NewStyleQueryHandler(repo)

	results, err := handler.SearchStyles(context.Background(), "test", 0)
	if err != nil {
		t.Fatalf("SearchStyles() error = %v", err)
	}

	// Should not error, limit defaults to 20
	_ = results
}
