package style

import (
	"errors"
	"time"
)

// Entity represents a prompt style in the domain
type Entity struct {
	ID          string
	Name        string
	Description string
	Prompt      string
	Category    string
	Tags        []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Validation errors
var (
	ErrEmptyName        = errors.New("style name cannot be empty")
	ErrEmptyPrompt      = errors.New("style prompt cannot be empty")
	ErrNameTooLong      = errors.New("style name cannot exceed 100 characters")
	ErrPromptTooLong    = errors.New("style prompt cannot exceed 2000 characters")
	ErrDescriptionTooLong = errors.New("style description cannot exceed 500 characters")
)

// New creates a new Style entity with validation
func New(name, description, prompt, category string, tags []string) (*Entity, error) {
	style := &Entity{
		Name:        name,
		Description: description,
		Prompt:      prompt,
		Category:    category,
		Tags:        tags,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := style.Validate(); err != nil {
		return nil, err
	}

	return style, nil
}

// Validate checks if the style entity is valid
func (s *Entity) Validate() error {
	if s.Name == "" {
		return ErrEmptyName
	}
	if len(s.Name) > 100 {
		return ErrNameTooLong
	}
	if s.Prompt == "" {
		return ErrEmptyPrompt
	}
	if len(s.Prompt) > 2000 {
		return ErrPromptTooLong
	}
	if len(s.Description) > 500 {
		return ErrDescriptionTooLong
	}
	return nil
}

// Update modifies the style entity
func (s *Entity) Update(name, description, prompt, category string, tags []string) error {
	s.Name = name
	s.Description = description
	s.Prompt = prompt
	s.Category = category
	s.Tags = tags
	s.UpdatedAt = time.Now()

	return s.Validate()
}
