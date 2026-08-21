package identity_test

import (
	"strings"
	"testing"

	"github.com/golangnigeria/curexal/internal/shared/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlphanumericCodeGeneration(t *testing.T) {
	t.Run("Generates 6 character uppercase code", func(t *testing.T) {
		code, err := crypto.GenerateAlphanumericCode(6)
		require.NoError(t, err)
		assert.Len(t, code, 6)
		assert.Equal(t, strings.ToUpper(code), code)

		// Ensure unambiguous charset (no 0, O, 1, I)
		for _, ch := range code {
			assert.NotContains(t, []rune{'0', 'O', '1', 'I'}, ch)
		}
	})

	t.Run("Generates distinct codes", func(t *testing.T) {
		code1, err1 := crypto.GenerateAlphanumericCode(6)
		code2, err2 := crypto.GenerateAlphanumericCode(6)
		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, code1, code2)
	})
}
