package realtime

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// RealtimePublisher is the vendor-agnostic interface for publishing live realtime events
// (e.g., patient queue updates, vitals triage alerts, consultation status changes).
type RealtimePublisher interface {
	Name() string
	Broadcast(ctx context.Context, channel string, event string, payload any) error
}

// Config holds settings for Realtime providers.
type Config struct {
	Provider string // "local", "memory", "noop"
}

// New creates a RealtimePublisher based on configuration.
func New(cfg Config) (RealtimePublisher, error) {
	provider := strings.ToLower(strings.TrimSpace(cfg.Provider))
	switch provider {
	case "local", "memory", "":
		return NewLocalRealtimePublisher(), nil
	case "noop":
		return &NoOpRealtimePublisher{}, nil
	default:
		return nil, fmt.Errorf("unsupported realtime provider %q (supported: local, memory, noop)", cfg.Provider)
	}
}

// LocalRealtimePublisher is an in-memory pub-sub implementation for local development and integration testing.
type LocalRealtimePublisher struct {
	mu          sync.RWMutex
	subscribers map[string][]chan any
	Events      []LocalEventRecord
}

type LocalEventRecord struct {
	Channel   string
	Event     string
	Payload   any
	Timestamp time.Time
}

func NewLocalRealtimePublisher() *LocalRealtimePublisher {
	return &LocalRealtimePublisher{
		subscribers: make(map[string][]chan any),
		Events:      make([]LocalEventRecord, 0),
	}
}

func (l *LocalRealtimePublisher) Name() string {
	return "local"
}

func (l *LocalRealtimePublisher) Broadcast(ctx context.Context, channel string, event string, payload any) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.Events = append(l.Events, LocalEventRecord{
		Channel:   channel,
		Event:     event,
		Payload:   payload,
		Timestamp: time.Now(),
	})

	if subs, ok := l.subscribers[channel]; ok {
		for _, ch := range subs {
			select {
			case ch <- payload:
			default:
			}
		}
	}

	return nil
}

// NoOpRealtimePublisher is a silent no-op publisher.
type NoOpRealtimePublisher struct{}

func (n *NoOpRealtimePublisher) Name() string {
	return "noop"
}

func (n *NoOpRealtimePublisher) Broadcast(ctx context.Context, channel string, event string, payload any) error {
	return nil
}
