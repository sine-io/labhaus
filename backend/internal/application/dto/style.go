package dto

import "time"

// StyleDTO represents a style for API responses
type StyleDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Prompt      string    `json:"prompt"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateStyleRequest represents a request to create a style
type CreateStyleRequest struct {
	Name        string   `json:"name" binding:"required,max=100"`
	Description string   `json:"description" binding:"max=500"`
	Prompt      string   `json:"prompt" binding:"required,max=2000"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
}

// UpdateStyleRequest represents a request to update a style
type UpdateStyleRequest struct {
	Name        string   `json:"name" binding:"required,max=100"`
	Description string   `json:"description" binding:"max=500"`
	Prompt      string   `json:"prompt" binding:"required,max=2000"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
}

// StyleFilterDTO represents filter options for style queries
type StyleFilterDTO struct {
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
}

// StyleListResponse represents a paginated list of styles
type StyleListResponse struct {
	Styles []StyleDTO `json:"styles"`
	Total  int        `json:"total"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}
