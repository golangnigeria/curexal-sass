package queue

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryQueueBroker(t *testing.T) {
	broker := NewMemoryQueueBroker(50)
	assert.Equal(t, "memory", broker.Name())

	var wg sync.WaitGroup
	var receivedPayload string

	broker.RegisterHandler("notifications.sms", func(ctx context.Context, job *Job) error {
		receivedPayload = string(job.Payload)
		wg.Done()
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := broker.Start(ctx)
	require.NoError(t, err)

	wg.Add(1)
	jobID, err := broker.Enqueue(ctx, "notifications.sms", []byte("Hello patient appointment confirmed"), 0)
	require.NoError(t, err)
	assert.NotEmpty(t, jobID)

	// Wait for async processing
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		assert.Equal(t, "Hello patient appointment confirmed", receivedPayload)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for job processing")
	}

	err = broker.Stop()
	require.NoError(t, err)
}
