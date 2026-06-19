package queue

import (
	"context"
	"encoding/json"
	"time"
)

// Queue represents an async task queue
type Queue interface {
	// Enqueue adds a task to the queue
	Enqueue(ctx context.Context, task *Task) error

	// Dequeue retrieves and removes a task from the queue (blocking)
	Dequeue(ctx context.Context, timeout time.Duration) (*Task, error)

	// StartWorker starts a worker goroutine to process tasks
	StartWorker(ctx context.Context, handler TaskHandler)

	// GetTaskStatus retrieves task status
	GetTaskStatus(ctx context.Context, taskID string) (*TaskStatus, error)

	// Close gracefully shuts down the queue
	Close() error
}

// Task represents a queued task
type Task struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"` // "workflow_execute", "image_generate"
	Payload    json.RawMessage `json:"payload"`
	RetryCount int             `json:"retry_count"`
	MaxRetries int             `json:"max_retries"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// TaskStatus represents the status of a task
type TaskStatus struct {
	TaskID     string    `json:"task_id"`
	State      string    `json:"state"` // "pending", "processing", "completed", "failed", "dead"
	RetryCount int       `json:"retry_count"`
	Error      string    `json:"error,omitempty"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TaskHandler is a function that processes a task
type TaskHandler func(ctx context.Context, task *Task) error

// TaskType constants
const (
	TaskTypeWorkflowExecute = "workflow_execute"
	TaskTypeImageGenerate   = "image_generate"
)

// Task states
const (
	TaskStatePending    = "pending"
	TaskStateProcessing = "processing"
	TaskStateCompleted  = "completed"
	TaskStateFailed     = "failed"
	TaskStateDead       = "dead" // exceeded max retries
)
