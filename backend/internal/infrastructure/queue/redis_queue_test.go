package queue_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/labhaus/backend/internal/infrastructure/queue"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRedisClient(t *testing.T) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // use test database
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	// Clean up test data
	client.FlushDB(ctx)

	return client
}

func TestRedisQueue_EnqueueDequeue(t *testing.T) {
	client := setupRedisClient(t)
	defer client.Close()

	q := queue.NewRedisQueue(client, "test")
	ctx := context.Background()

	// Create test task
	payload, _ := json.Marshal(map[string]string{"key": "value"})
	task := &queue.Task{
		Type:       queue.TaskTypeWorkflowExecute,
		Payload:    payload,
		MaxRetries: 3,
	}

	// Enqueue
	err := q.Enqueue(ctx, task)
	require.NoError(t, err)
	assert.NotEmpty(t, task.ID)

	// Dequeue
	dequeuedTask, err := q.Dequeue(ctx, 1*time.Second)
	require.NoError(t, err)
	require.NotNil(t, dequeuedTask)

	assert.Equal(t, task.ID, dequeuedTask.ID)
	assert.Equal(t, task.Type, dequeuedTask.Type)
	assert.JSONEq(t, string(task.Payload), string(dequeuedTask.Payload))
}

func TestRedisQueue_Worker(t *testing.T) {
	client := setupRedisClient(t)
	defer client.Close()

	q := queue.NewRedisQueue(client, "test_worker")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	processed := make(chan string, 1)

	// Start worker
	handler := func(ctx context.Context, task *queue.Task) error {
		processed <- task.ID
		return nil
	}
	q.StartWorker(ctx, handler)

	// Enqueue task
	task := &queue.Task{
		Type:    queue.TaskTypeWorkflowExecute,
		Payload: json.RawMessage(`{"test": true}`),
	}
	err := q.Enqueue(ctx, task)
	require.NoError(t, err)

	// Wait for processing
	select {
	case taskID := <-processed:
		assert.Equal(t, task.ID, taskID)
	case <-time.After(3 * time.Second):
		t.Fatal("Task was not processed")
	}

	// Check task status
	status, err := q.GetTaskStatus(ctx, task.ID)
	require.NoError(t, err)
	assert.Equal(t, queue.TaskStateCompleted, status.State)
}

func TestRedisQueue_Retry(t *testing.T) {
	client := setupRedisClient(t)
	defer client.Close()

	q := queue.NewRedisQueue(client, "test_retry")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	attempts := 0
	processed := make(chan int, 5)

	// Handler that fails first 2 times
	handler := func(ctx context.Context, task *queue.Task) error {
		attempts++
		processed <- attempts
		if attempts < 3 {
			return assert.AnError // fail
		}
		return nil // succeed on 3rd attempt
	}
	q.StartWorker(ctx, handler)

	// Enqueue task
	task := &queue.Task{
		Type:       queue.TaskTypeWorkflowExecute,
		Payload:    json.RawMessage(`{"test": true}`),
		MaxRetries: 3,
	}
	err := q.Enqueue(ctx, task)
	require.NoError(t, err)

	// Wait for all attempts
	for i := 0; i < 3; i++ {
		select {
		case <-processed:
			// received attempt
		case <-time.After(5 * time.Second):
			t.Fatalf("Timeout waiting for attempt %d", i+1)
		}
	}

	assert.Equal(t, 3, attempts)

	// Check final status
	time.Sleep(500 * time.Millisecond)
	status, err := q.GetTaskStatus(ctx, task.ID)
	require.NoError(t, err)
	assert.Equal(t, queue.TaskStateCompleted, status.State)
}

func TestRedisQueue_DeadLetterQueue(t *testing.T) {
	client := setupRedisClient(t)
	defer client.Close()

	q := queue.NewRedisQueue(client, "test_dlq")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Handler that always fails
	handler := func(ctx context.Context, task *queue.Task) error {
		return assert.AnError
	}
	q.StartWorker(ctx, handler)

	// Enqueue task with max 2 retries
	task := &queue.Task{
		Type:       queue.TaskTypeWorkflowExecute,
		Payload:    json.RawMessage(`{"test": true}`),
		MaxRetries: 2,
	}
	err := q.Enqueue(ctx, task)
	require.NoError(t, err)

	// Wait for all retries to exhaust
	time.Sleep(3 * time.Second)

	// Check task moved to dead state
	status, err := q.GetTaskStatus(ctx, task.ID)
	require.NoError(t, err)
	assert.Equal(t, queue.TaskStateDead, status.State)
	assert.NotEmpty(t, status.Error)
}
