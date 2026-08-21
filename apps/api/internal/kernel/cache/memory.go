package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type cacheItem struct {
	value      interface{}
	expiration int64
}

func (item cacheItem) isExpired() bool {
	if item.expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.expiration
}

// MemoryCache provides a fast, thread-safe, zero-dependency in-memory TTL cache.
type MemoryCache struct {
	mu         sync.RWMutex
	items      map[string]cacheItem
	defaultTTL time.Duration
}

// NewMemoryCache creates a new in-memory cache instance.
func NewMemoryCache(defaultTTL time.Duration) *MemoryCache {
	if defaultTTL <= 0 {
		defaultTTL = 5 * time.Minute
	}
	c := &MemoryCache{
		items:      make(map[string]cacheItem),
		defaultTTL: defaultTTL,
	}

	// Background cleanup goroutine for expired items every 2 minutes
	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		for range ticker.C {
			c.cleanupExpired()
		}
	}()

	return c
}

func (c *MemoryCache) cleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now().UnixNano()
	for k, v := range c.items {
		if v.expiration > 0 && now > v.expiration {
			delete(c.items, k)
		}
	}
}

func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found || item.isExpired() {
		return nil, false
	}
	return item.value, true
}

func (c *MemoryCache) GetString(ctx context.Context, key string) (string, bool) {
	val, ok := c.Get(ctx, key)
	if !ok {
		return "", false
	}
	str, ok := val.(string)
	if !ok {
		return fmt.Sprintf("%v", val), true
	}
	return str, true
}

func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = c.defaultTTL
	}

	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	c.mu.Lock()
	c.items[key] = cacheItem{
		value:      value,
		expiration: exp,
	}
	c.mu.Unlock()
	return nil
}

func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
	return nil
}

func (c *MemoryCache) Exists(ctx context.Context, key string) bool {
	_, ok := c.Get(ctx, key)
	return ok
}

func (c *MemoryCache) Flush(ctx context.Context) error {
	c.mu.Lock()
	c.items = make(map[string]cacheItem)
	c.mu.Unlock()
	return nil
}
