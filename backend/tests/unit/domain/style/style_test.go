package style_test

import (
	"testing"

	"github.com/labhaus/backend/internal/domain/style"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		styleName   string
		description string
		prompt      string
		category    string
		tags        []string
		wantErr     error
	}{
		{
			name:        "valid style",
			styleName:   "Anime Style",
			description: "Japanese anime art style",
			prompt:      "anime, vibrant colors, expressive eyes",
			category:    "Art",
			tags:        []string{"anime", "japanese"},
			wantErr:     nil,
		},
		{
			name:        "empty name",
			styleName:   "",
			description: "Test",
			prompt:      "test prompt",
			category:    "Test",
			tags:        []string{},
			wantErr:     style.ErrEmptyName,
		},
		{
			name:        "empty prompt",
			styleName:   "Test Style",
			description: "Test",
			prompt:      "",
			category:    "Test",
			tags:        []string{},
			wantErr:     style.ErrEmptyPrompt,
		},
		{
			name:        "name too long",
			styleName:   string(make([]byte, 101)),
			description: "Test",
			prompt:      "test prompt",
			category:    "Test",
			tags:        []string{},
			wantErr:     style.ErrNameTooLong,
		},
		{
			name:        "prompt too long",
			styleName:   "Test Style",
			description: "Test",
			prompt:      string(make([]byte, 2001)),
			category:    "Test",
			tags:        []string{},
			wantErr:     style.ErrPromptTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := style.New(tt.styleName, tt.description, tt.prompt, tt.category, tt.tags)
			if err != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && s == nil {
				t.Error("New() returned nil style with no error")
			}
			if err == nil {
				if s.Name != tt.styleName {
					t.Errorf("Name = %v, want %v", s.Name, tt.styleName)
				}
				if s.CreatedAt.IsZero() {
					t.Error("CreatedAt should not be zero")
				}
				if s.UpdatedAt.IsZero() {
					t.Error("UpdatedAt should not be zero")
				}
			}
		})
	}
}

func TestEntity_Update(t *testing.T) {
	s, err := style.New("Original", "Original desc", "Original prompt", "Cat", []string{"tag1"})
	if err != nil {
		t.Fatalf("Failed to create style: %v", err)
	}

	originalUpdatedAt := s.UpdatedAt

	err = s.Update("Updated", "Updated desc", "Updated prompt", "NewCat", []string{"tag2"})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if s.Name != "Updated" {
		t.Errorf("Name = %v, want Updated", s.Name)
	}
	if s.Description != "Updated desc" {
		t.Errorf("Description = %v, want Updated desc", s.Description)
	}
	if s.UpdatedAt == originalUpdatedAt {
		t.Error("UpdatedAt should be updated")
	}

	// Test invalid update
	err = s.Update("", "desc", "prompt", "cat", []string{})
	if err != style.ErrEmptyName {
		t.Errorf("Update() with empty name should return ErrEmptyName, got %v", err)
	}
}

func TestEntity_Validate(t *testing.T) {
	validStyle, _ := style.New("Test", "Test desc", "Test prompt", "Cat", []string{})

	if err := validStyle.Validate(); err != nil {
		t.Errorf("Validate() should pass for valid style, got %v", err)
	}

	// Manually create invalid style
	invalidStyle := &style.Entity{
		Name:   "",
		Prompt: "test",
	}

	if err := invalidStyle.Validate(); err != style.ErrEmptyName {
		t.Errorf("Validate() should return ErrEmptyName, got %v", err)
	}
}
