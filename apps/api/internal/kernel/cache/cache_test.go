package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/cache"
)

func TestMemoryCache(t *testing.T) {
	c := cache.New(cache.Config{
		Provider:   "memory",
		DefaultTTL: 100 * time.Millisecond,
	})

	ctx := context.Background()

	// Set & Get
	err := c.Set(ctx, "sample_key", "sample_value", 500*time.Millisecond)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, found := c.GetString(ctx, "sample_key")
	if !found || val != "sample_value" {
		t.Fatalf("expected sample_value, got %s (found: %v)", val, found)
	}

	// Exists
	if !c.Exists(ctx, "sample_key") {
		t.Fatalf("expected key to exist")
	}

	// Delete
	_ = c.Delete(ctx, "sample_key")
	if c.Exists(ctx, "sample_key") {
		t.Fatalf("expected key to be deleted")
	}

	// Expiration
	_ = c.Set(ctx, "expiring_key", "will_expire", 50*time.Millisecond)
	time.Sleep(70 * time.Millisecond)
	_, foundExp := c.Get(ctx, "expiring_key")
	if foundExp {
		t.Fatalf("expected expired key to not be found")
	}
}
