package realtime

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRealtimePublisher_Factory(t *testing.T) {
	t.Run("Create Local Publisher", func(t *testing.T) {
		pub, err := New(Config{Provider: "local"})
		require.NoError(t, err)
		assert.Equal(t, "local", pub.Name())

		err = pub.Broadcast(context.Background(), "clinic:org_123:queue", "appointment.checked_in", map[string]string{
			"patient_id": "pt_123",
			"status":     "in_triage",
		})
		require.NoError(t, err)

		localPub := pub.(*LocalRealtimePublisher)
		assert.Len(t, localPub.Events, 1)
		assert.Equal(t, "clinic:org_123:queue", localPub.Events[0].Channel)
		assert.Equal(t, "appointment.checked_in", localPub.Events[0].Event)
	})

	t.Run("Create NoOp Publisher", func(t *testing.T) {
		pub, err := New(Config{Provider: "noop"})
		require.NoError(t, err)
		assert.Equal(t, "noop", pub.Name())

		err = pub.Broadcast(context.Background(), "channel", "event", nil)
		require.NoError(t, err)
	})
}
