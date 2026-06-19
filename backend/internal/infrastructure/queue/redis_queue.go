package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RedisQueue implements Queue using Redis
type RedisQueue struct {
	client     *redis.Client
	queueKey   string
	dlqKey     string // dead letter queue
	taskPrefix string
	workers    int
	stopCh     chan struct{}
}

// NewRedisQueue creates a new Redis-backed queue
func NewRedisQueue(client *redis.Client, queueName string) *RedisQueue {
	return &RedisQueue{
		client:     client,
		queueKey:   fmt.Sprintf("queue:%s", queueName),
		dlqKey:     fmt.Sprintf("queue:%s:dlq", queueName),
		taskPrefix: fmt.Sprintf("task:%s:", queueName),
		workers:    1,
		stopCh:     make(chan struct{}),
	}
}

// Enqueue adds a task to the queue
func (q *RedisQueue) Enqueue(ctx context.Context, task *Task) error {
	// Generate ID if not set
	if task.ID == "" {
		task.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	if task.CreatedAt.IsZero() {
		task.CreatedAt = now
	}
	task.UpdatedAt = now

	// Set default max retries
	if task.MaxRetries == 0 {
		task.MaxRetries = 3
	}

	// Serialize task
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	// Store task metadata in hash
	taskKey := q.taskKey(task.ID)
	if err := q.client.HSet(ctx, taskKey, map[string]interface{}{
		"data":   string(taskJSON),
		"state":  TaskStatePending,
		"enqueued_at": now.Unix(),
	}).Err(); err != nil {
		return fmt.Errorf("failed to store task metadata: %w", err)
	}

	// Set expiration (24 hours)
	q.client.Expire(ctx, taskKey, 24*time.Hour)

	// Push task ID to queue
	if err := q.client.LPush(ctx, q.queueKey, task.ID).Err(); err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	return nil
}

// Dequeue retrieves and removes a task from the queue (blocking)
func (q *RedisQueue) Dequeue(ctx context.Context, timeout time.Duration) (*Task, error) {
	// BRPOP with timeout
	result, err := q.client.BRPop(ctx, timeout, q.queueKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // timeout, no task available
		}
		return nil, fmt.Errorf("failed to dequeue: %w", err)
	}

	// result[0] is the key, result[1] is the value (task ID)
	if len(result) < 2 {
		return nil, errors.New("invalid dequeue result")
	}

	taskID := result[1]
	taskKey := q.taskKey(taskID)

	// Get task data
	taskData, err := q.client.HGet(ctx, taskKey, "data").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("task %s not found", taskID)
		}
		return nil, fmt.Errorf("failed to get task data: %w", err)
	}

	// Deserialize task
	var task Task
	if err := json.Unmarshal([]byte(taskData), &task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task: %w", err)
	}

	// Update task state to processing
	q.client.HSet(ctx, taskKey, "state", TaskStateProcessing)
	q.client.HSet(ctx, taskKey, "processing_at", time.Now().Unix())

	return &task, nil
}

// StartWorker starts a worker goroutine to process tasks
func (q *RedisQueue) StartWorker(ctx context.Context, handler TaskHandler) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-q.stopCh:
				return
			default:
				// Dequeue with 5 second timeout
				task, err := q.Dequeue(ctx, 5*time.Second)
				if err != nil {
					// Log error but continue
					fmt.Printf("Dequeue error: %v\n", err)
					time.Sleep(time.Second)
					continue
				}

				if task == nil {
					// No task available, continue
					continue
				}

				// Process task
				if err := q.processTask(ctx, task, handler); err != nil {
					fmt.Printf("Task processing error: %v\n", err)
				}
			}
		}
	}()
}

// processTask handles task execution and retry logic
func (q *RedisQueue) processTask(ctx context.Context, task *Task, handler TaskHandler) error {
	taskKey := q.taskKey(task.ID)

	// Execute handler
	err := handler(ctx, task)

	if err != nil {
		// Task failed
		task.RetryCount++
		task.UpdatedAt = time.Now()

		if task.RetryCount >= task.MaxRetries {
			// Max retries exceeded, move to dead letter queue
			q.client.HSet(ctx, taskKey, "state", TaskStateDead)
			q.client.HSet(ctx, taskKey, "error", err.Error())
			q.client.HSet(ctx, taskKey, "failed_at", time.Now().Unix())

			// Move to DLQ
			taskJSON, _ := json.Marshal(task)
			q.client.LPush(ctx, q.dlqKey, string(taskJSON))

			return fmt.Errorf("task %s exceeded max retries: %w", task.ID, err)
		}

		// Retry: re-enqueue task
		q.client.HSet(ctx, taskKey, "state", TaskStatePending)
		q.client.HSet(ctx, taskKey, "retry_count", task.RetryCount)
		q.client.HSet(ctx, taskKey, "error", err.Error())

		// Update task data
		taskJSON, _ := json.Marshal(task)
		q.client.HSet(ctx, taskKey, "data", string(taskJSON))

		// Re-enqueue
		q.client.LPush(ctx, q.queueKey, task.ID)

		return fmt.Errorf("task %s failed, retrying (%d/%d): %w", task.ID, task.RetryCount, task.MaxRetries, err)
	}

	// Task succeeded
	q.client.HSet(ctx, taskKey, "state", TaskStateCompleted)
	q.client.HSet(ctx, taskKey, "completed_at", time.Now().Unix())

	return nil
}

// GetTaskStatus retrieves task status
func (q *RedisQueue) GetTaskStatus(ctx context.Context, taskID string) (*TaskStatus, error) {
	taskKey := q.taskKey(taskID)

	// Get task metadata
	result, err := q.client.HGetAll(ctx, taskKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get task status: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("task %s not found", taskID)
	}

	status := &TaskStatus{
		TaskID: taskID,
		State:  result["state"],
		Error:  result["error"],
	}

	// Parse retry count
	if retryCountStr, ok := result["retry_count"]; ok {
		fmt.Sscanf(retryCountStr, "%d", &status.RetryCount)
	}

	// Parse updated timestamp
	if updatedAtStr, ok := result["completed_at"]; ok {
		var timestamp int64
		fmt.Sscanf(updatedAtStr, "%d", &timestamp)
		status.UpdatedAt = time.Unix(timestamp, 0)
	} else if updatedAtStr, ok := result["failed_at"]; ok {
		var timestamp int64
		fmt.Sscanf(updatedAtStr, "%d", &timestamp)
		status.UpdatedAt = time.Unix(timestamp, 0)
	} else if updatedAtStr, ok := result["processing_at"]; ok {
		var timestamp int64
		fmt.Sscanf(updatedAtStr, "%d", &timestamp)
		status.UpdatedAt = time.Unix(timestamp, 0)
	} else if updatedAtStr, ok := result["enqueued_at"]; ok {
		var timestamp int64
		fmt.Sscanf(updatedAtStr, "%d", &timestamp)
		status.UpdatedAt = time.Unix(timestamp, 0)
	}

	return status, nil
}

// Close gracefully shuts down the queue
func (q *RedisQueue) Close() error {
	close(q.stopCh)
	return nil
}

// taskKey generates Redis key for task metadata
func (q *RedisQueue) taskKey(taskID string) string {
	return q.taskPrefix + taskID
}
