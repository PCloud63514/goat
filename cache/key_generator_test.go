package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_KeyGenerator(t *testing.T) {
	gen := &keyGenerator{}

	t.Run("should generate key with multiple values", func(t *testing.T) {
		key := gen.Generate("user", 123, true)
		assert.Equal(t, "user-123-true", key)
	})

	t.Run("should generate key with single value", func(t *testing.T) {
		key := gen.Generate("single")
		assert.Equal(t, "single", key)
	})

	t.Run("should panic when no values provided", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic on empty params, but got none")
			}
		}()
		_ = gen.Generate()
	})

	t.Run("should convert various types to string", func(t *testing.T) {
		key := gen.Generate("str", 42, 3.14, false)
		assert.Equal(t, "str-42-3.14-false", key)
	})
}
