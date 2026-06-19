package queue

import (
	"context"
	"encoding/json"
	"fmt"
)

// WorkflowTaskPayload represents the payload for workflow execution tasks
type WorkflowTaskPayload struct {
	WorkflowID string `json:"workflow_id"`
}

// NewWorkflowTask creates a new workflow execution task
func NewWorkflowTask(workflowID string) (*Task, error) {
	payload, err := json.Marshal(WorkflowTaskPayload{
		WorkflowID: workflowID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return &Task{
		Type:       TaskTypeWorkflowExecute,
		Payload:    payload,
		MaxRetries: 3,
	}, nil
}

// WorkflowTaskHandler creates a handler for workflow execution tasks
func WorkflowTaskHandler() TaskHandler {
	return func(ctx context.Context, task *Task) error {
		// Parse payload
		var payload WorkflowTaskPayload
		if err := json.Unmarshal(task.Payload, &payload); err != nil {
			return fmt.Errorf("invalid payload: %w", err)
		}

		// TODO: Implement actual workflow execution logic
		// For now, just log
		fmt.Printf("Processing workflow task: workflow_id=%s, attempt=%d/%d\n",
			payload.WorkflowID, task.RetryCount+1, task.MaxRetries)

		// Simulate work
		// In production, this would:
		// 1. Load workflow from database
		// 2. Execute workflow steps
		// 3. Update workflow status
		// 4. Store results

		return nil
	}
}
