package cache

import (
	"context"
	"time"
)

// CacheService defines the interface for ephemeral caching operations.
// It is intended ONLY for non-authoritative optimizations (e.g. bootstrap caching, catalog queries).
// All critical data (sessions, billing, medical records) must reside durably in PostgreSQL.
type CacheService interface {
	Get(ctx context.Context, key string) (interface{}, bool)
	GetString(ctx context.Context, key string) (string, bool)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
	Flush(ctx context.Context) error
}

// Config represents cache initialization configuration.
type Config struct {
	Provider     string        // "memory" or "redis"
	RedisAddress string        // Redis address if available (optional)
	DefaultTTL   time.Duration // Default cache duration
}

// New initializes the appropriate CacheService implementation.
func New(cfg Config) CacheService {
	// For Dockerless production and dev baselines, the in-memory cache is fully self-contained.
	return NewMemoryCache(cfg.DefaultTTL)
}
