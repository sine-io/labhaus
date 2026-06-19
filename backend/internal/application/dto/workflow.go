package dto

import "time"

// WorkflowDTO represents a workflow for API responses
type WorkflowDTO struct {
	ID        string        `json:"id"`
	UserID    string        `json:"user_id"`
	StyleID   string        `json:"style_id"`
	State     string        `json:"state"`
	Config    WorkflowConfig `json:"config"`
	Result    *WorkflowResult `json:"result,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// WorkflowConfig represents workflow configuration
type WorkflowConfig struct {
	ImageCount int   `json:"image_count" binding:"required,min=1,max=10"`
	Width      int   `json:"width" binding:"required,min=1"`
	Height     int   `json:"height" binding:"required,min=1"`
	Steps      int   `json:"steps" binding:"min=1,max=100"`
	Seed       int64 `json:"seed"`
}

// WorkflowResult represents workflow execution result
type WorkflowResult struct {
	ImageURLs []string `json:"image_urls"`
	Duration  int64    `json:"duration_ms"` // milliseconds
	Error     string   `json:"error,omitempty"`
}

// CreateWorkflowRequest represents a request to create a workflow
type CreateWorkflowRequest struct {
	StyleID string         `json:"style_id" binding:"required"`
	Config  WorkflowConfig `json:"config" binding:"required"`
}

// UpdateWorkflowStateRequest represents a request to update workflow state
type UpdateWorkflowStateRequest struct {
	State string `json:"state" binding:"required"`
}

// WorkflowFilterDTO represents filter options for workflow queries
type WorkflowFilterDTO struct {
	State  *string `json:"state"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
}

// WorkflowListResponse represents a paginated list of workflows
type WorkflowListResponse struct {
	Workflows []WorkflowDTO `json:"workflows"`
	Total     int           `json:"total"`
	Limit     int           `json:"limit"`
	Offset    int           `json:"offset"`
}
