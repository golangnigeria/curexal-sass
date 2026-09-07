package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Job represents a background unit of work.
type Job struct {
	ID          string            `json:"id"`
	Queue       string            `json:"queue"`
	Payload     []byte            `json:"payload"`
	Attempts    int               `json:"attempts"`
	MaxAttempts int               `json:"max_attempts"`
	CreatedAt   time.Time         `json:"created_at"`
	ScheduledAt time.Time         `json:"scheduled_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// HandlerFunc processes a queued job.
type HandlerFunc func(ctx context.Context, job *Job) error

// QueueBroker is the canonical vendor-agnostic interface for queuing background tasks
// (e.g., patient SMS reminders, invoice generation, outbox dispatch, webhook delivery).
type QueueBroker interface {
	Name() string
	Enqueue(ctx context.Context, queue string, payload []byte, delay time.Duration) (string, error)
	RegisterHandler(queue string, handler HandlerFunc)
	Start(ctx context.Context) error
	Stop() error
}

// MemoryQueueBroker is an in-memory queue broker for testing and low-volume local development.
type MemoryQueueBroker struct {
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
	jobs     chan *Job
	stopChan chan struct{}
}

func NewMemoryQueueBroker(bufferSize int) *MemoryQueueBroker {
	if bufferSize <= 0 {
		bufferSize = 100
	}
	return &MemoryQueueBroker{
		handlers: make(map[string]HandlerFunc),
		jobs:     make(chan *Job, bufferSize),
		stopChan: make(chan struct{}),
	}
}

func (m *MemoryQueueBroker) Name() string {
	return "memory"
}

func (m *MemoryQueueBroker) RegisterHandler(queue string, handler HandlerFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[queue] = handler
}

func (m *MemoryQueueBroker) Enqueue(ctx context.Context, queue string, payload []byte, delay time.Duration) (string, error) {
	jobID := fmt.Sprintf("job_mem_%d", time.Now().UnixNano())
	job := &Job{
		ID:          jobID,
		Queue:       queue,
		Payload:     payload,
		Attempts:    0,
		MaxAttempts: 3,
		CreatedAt:   time.Now(),
		ScheduledAt: time.Now().Add(delay),
	}

	if delay > 0 {
		go func() {
			select {
			case <-time.After(delay):
				m.jobs <- job
			case <-m.stopChan:
			}
		}()
	} else {
		select {
		case m.jobs <- job:
		default:
			return "", errors.New("queue buffer is full")
		}
	}

	return jobID, nil
}

func (m *MemoryQueueBroker) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-m.stopChan:
				return
			case <-ctx.Done():
				return
			case job := <-m.jobs:
				m.mu.RLock()
				handler, ok := m.handlers[job.Queue]
				m.mu.RUnlock()

				if ok && handler != nil {
					_ = handler(context.Background(), job)
				}
			}
		}
	}()
	return nil
}

func (m *MemoryQueueBroker) Stop() error {
	close(m.stopChan)
	return nil
}

// PGMQBroker implements QueueBroker utilizing PostgreSQL / Supabase pgmq extension.
type PGMQBroker struct {
	pool     *pgxpool.Pool
	handlers map[string]HandlerFunc
	stopChan chan struct{}
	mu       sync.RWMutex
}

func NewPGMQBroker(pool *pgxpool.Pool) *PGMQBroker {
	return &PGMQBroker{
		pool:     pool,
		handlers: make(map[string]HandlerFunc),
		stopChan: make(chan struct{}),
	}
}

func (p *PGMQBroker) Name() string {
	return "pgmq"
}

func (p *PGMQBroker) RegisterHandler(queue string, handler HandlerFunc) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers[queue] = handler
}

func (p *PGMQBroker) Enqueue(ctx context.Context, queue string, payload []byte, delay time.Duration) (string, error) {
	if p.pool == nil {
		return "", errors.New("database connection pool not initialized")
	}

	delaySeconds := int(delay.Seconds())
	var msgID int64
	query := "SELECT pgmq.send($1, $2::jsonb, $3);"
	err := p.pool.QueryRow(ctx, query, queue, string(payload), delaySeconds).Scan(&msgID)
	if err != nil {
		// Fallback for raw JSON payload if pgmq function signature expects jsonb
		return "", fmt.Errorf("pgmq enqueue failed: %w", err)
	}

	return fmt.Sprintf("%d", msgID), nil
}

func (p *PGMQBroker) Start(ctx context.Context) error {
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case <-p.stopChan:
				ticker.Stop()
				return
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				p.pollQueues(ctx)
			}
		}
	}()
	return nil
}

func (p *PGMQBroker) pollQueues(ctx context.Context) {
	p.mu.RLock()
	queues := make([]string, 0, len(p.handlers))
	for q := range p.handlers {
		queues = append(queues, q)
	}
	p.mu.RUnlock()

	for _, queueName := range queues {
		p.processBatch(ctx, queueName)
	}
}

func (p *PGMQBroker) processBatch(ctx context.Context, queueName string) {
	if p.pool == nil {
		return
	}

	query := "SELECT msg_id, read_ct, message::text FROM pgmq.read($1, 30, 5);"
	rows, err := p.pool.Query(ctx, query, queueName)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var msgID int64
		var readCt int
		var messageStr string
		if err := rows.Scan(&msgID, &readCt, &messageStr); err != nil {
			continue
		}

		p.mu.RLock()
		handler := p.handlers[queueName]
		p.mu.RUnlock()

		if handler != nil {
			job := &Job{
				ID:          fmt.Sprintf("%d", msgID),
				Queue:       queueName,
				Payload:     []byte(messageStr),
				Attempts:    readCt,
				MaxAttempts: 5,
				CreatedAt:   time.Now(),
			}

			if err := handler(ctx, job); err == nil {
				// Delete successfully processed message from pgmq
				_, _ = p.pool.Exec(ctx, "SELECT pgmq.delete($1, $2);", queueName, msgID)
			}
		}
	}
}

func (p *PGMQBroker) Stop() error {
	close(p.stopChan)
	return nil
}

// EnqueueTyped is a convenience helper for enqueuing any JSON-serializable struct.
func EnqueueTyped[T any](ctx context.Context, broker QueueBroker, queue string, payload T, delay time.Duration) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to serialize payload: %w", err)
	}
	return broker.Enqueue(ctx, queue, data, delay)
}
