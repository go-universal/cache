package cache_test

import (
	"testing"

	"github.com/go-universal/cache"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisQueue(t *testing.T) {
	queue := cache.NewRedisQueue("test-queue", redis.NewClient(&redis.Options{}))

	t.Run("Push and Length", func(t *testing.T) {
		err := queue.Clear()
		require.NoError(t, err, "Failed to clear queue")

		err = queue.Push("item1")
		require.NoError(t, err, "Failed to push item to queue")

		length, err := queue.Length()
		require.NoError(t, err, "Failed to get queue length")
		assert.Equal(t, int64(1), length, "Queue length should be 1")
	})

	t.Run("Pull", func(t *testing.T) {
		err := queue.Clear()
		require.NoError(t, err, "Failed to clear queue")

		err = queue.Push("item2")
		require.NoError(t, err, "Failed to push item to queue")

		item, err := queue.Pull()
		require.NoError(t, err, "Failed to pull item from queue")
		assert.Equal(t, "item2", item, "Pulled item should match pushed item")
	})

	t.Run("Pop", func(t *testing.T) {
		err := queue.Clear()
		require.NoError(t, err, "Failed to clear queue")

		err = queue.Push("item3")
		require.NoError(t, err, "Failed to push item to queue")

		item, err := queue.Pop()
		require.NoError(t, err, "Failed to pop item from queue")
		assert.Equal(t, "item3", item, "Popped item should match pushed item")
	})

	t.Run("Cast", func(t *testing.T) {
		err := queue.Clear()
		require.NoError(t, err, "Failed to clear queue")

		err = queue.Push("item4")
		require.NoError(t, err, "Failed to push item to queue")

		caster, err := queue.Cast()
		require.NoError(t, err, "Failed to cast item from queue")
		assert.Equal(t, "item4", caster.StringSafe(""), "Casted item should match pushed item")
	})
}
