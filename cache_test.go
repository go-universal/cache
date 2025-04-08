package cache_test

import (
	"testing"
	"time"

	"github.com/go-universal/cache"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryCache(t *testing.T) {
	memCache := cache.NewMemoryCache()

	t.Run("Put and Get", func(t *testing.T) {
		key := "testKey"
		value := "testValue"
		ttl := 5 * time.Second

		err := memCache.Put(key, value, &ttl)
		require.NoError(t, err)

		retrievedValue, err := memCache.Get(key)
		require.NoError(t, err)
		assert.Equal(t, value, retrievedValue)
	})

	t.Run("Update existing key", func(t *testing.T) {
		key := "testKey"
		newValue := "newValue"

		exists, err := memCache.Update(key, newValue)
		require.NoError(t, err)
		assert.True(t, exists)

		retrievedValue, err := memCache.Get(key)
		require.NoError(t, err)
		assert.Equal(t, newValue, retrievedValue)
	})

	t.Run("PutOrUpdate with TTL", func(t *testing.T) {
		key := "testKey"
		value := "overriddenValue"
		ttl := 10 * time.Second

		err := memCache.PutOrUpdate(key, value, &ttl)
		require.NoError(t, err)

		retrievedValue, err := memCache.Get(key)
		require.NoError(t, err)
		assert.Equal(t, value, retrievedValue)
	})

	t.Run("Exists", func(t *testing.T) {
		key := "testKey"

		exists, err := memCache.Exists(key)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Forget", func(t *testing.T) {
		key := "testKey"

		err := memCache.Forget(key)
		require.NoError(t, err)

		exists, err := memCache.Exists(key)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Increment and Decrement", func(t *testing.T) {
		key := "counter"
		initialValue := int64(10)
		err := memCache.Put(key, initialValue, nil)
		require.NoError(t, err)

		incremented, err := memCache.Increment(key, 5)
		require.NoError(t, err)
		assert.True(t, incremented)

		value, err := memCache.Get(key)
		require.NoError(t, err)
		assert.Equal(t, int64(15), value)

		decremented, err := memCache.Decrement(key, 3)
		require.NoError(t, err)
		assert.True(t, decremented)

		value, err = memCache.Get(key)
		require.NoError(t, err)
		assert.Equal(t, int64(12), value)
	})

	t.Run("Increment and Decrement Float", func(t *testing.T) {
		key := "counter"
		initialValue := float64(10.3)
		err := memCache.Put(key, initialValue, nil)
		require.NoError(t, err)

		incremented, err := memCache.IncrementFloat(key, 5)
		require.NoError(t, err)
		assert.True(t, incremented)

		value, err := memCache.Get(key)
		require.NoError(t, err)
		assert.Equal(t, float64(15.3), value)

		decremented, err := memCache.DecrementFloat(key, 3)
		require.NoError(t, err)
		assert.True(t, decremented)

		value, err = memCache.Get(key)
		require.NoError(t, err)
		assert.Equal(t, float64(12.3), value)
	})

	t.Run("TTL", func(t *testing.T) {
		key := "ttlKey"
		value := "ttlValue"
		ttl := 2 * time.Second

		err := memCache.Put(key, value, &ttl)
		require.NoError(t, err)

		retrievedTTL, err := memCache.TTL(key)
		require.NoError(t, err)
		assert.LessOrEqual(t, retrievedTTL.Seconds(), ttl.Seconds())
	})
}

func TestRedisCache(t *testing.T) {
	redisCache := cache.NewRedisCache("test", redis.NewClient(&redis.Options{}))

	t.Run("Put and Get", func(t *testing.T) {
		key := "redisTestKey"
		value := "redisTestValue"
		ttl := 5 * time.Second

		err := redisCache.Put(key, value, &ttl)
		require.NoError(t, err)

		retrievedValue, err := redisCache.Get(key)
		require.NoError(t, err)
		assert.Equal(t, value, retrievedValue)
	})

	t.Run("Update existing key", func(t *testing.T) {
		key := "redisTestKey"
		newValue := "redisNewValue"

		exists, err := redisCache.Update(key, newValue)
		require.NoError(t, err)
		assert.True(t, exists)

		retrievedValue, err := redisCache.Get(key)
		require.NoError(t, err)
		assert.Equal(t, newValue, retrievedValue)
	})

	t.Run("PutOrUpdate with TTL", func(t *testing.T) {
		key := "redisTestKey"
		value := "redisOverriddenValue"
		ttl := 10 * time.Second

		err := redisCache.PutOrUpdate(key, value, &ttl)
		require.NoError(t, err)

		retrievedValue, err := redisCache.Get(key)
		require.NoError(t, err)
		assert.Equal(t, value, retrievedValue)
	})

	t.Run("Exists", func(t *testing.T) {
		key := "redisTestKey"

		exists, err := redisCache.Exists(key)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Forget", func(t *testing.T) {
		key := "redisTestKey"

		err := redisCache.Forget(key)
		require.NoError(t, err)

		exists, err := redisCache.Exists(key)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Increment and Decrement", func(t *testing.T) {
		key := "redisCounter"
		initialValue := int64(10)
		err := redisCache.Put(key, initialValue, nil)
		require.NoError(t, err)

		incremented, err := redisCache.Increment(key, 5)
		require.NoError(t, err)
		assert.True(t, incremented)

		value, err := redisCache.Cast(key)
		require.NoError(t, err)
		assert.Equal(t, int64(15), value.Int64Safe(0))

		decremented, err := redisCache.Decrement(key, 3)
		require.NoError(t, err)
		assert.True(t, decremented)

		value, err = redisCache.Cast(key)
		require.NoError(t, err)
		assert.Equal(t, int64(12), value.Int64Safe(0))
	})

	t.Run("Increment and Decrement Float", func(t *testing.T) {
		key := "counter"
		initialValue := float64(10.3)
		err := redisCache.Put(key, initialValue, nil)
		require.NoError(t, err)

		incremented, err := redisCache.IncrementFloat(key, 5)
		require.NoError(t, err)
		assert.True(t, incremented)

		value, err := redisCache.Cast(key)
		require.NoError(t, err)
		assert.Equal(t, float64(15.3), value.Float64Safe(0))

		decremented, err := redisCache.DecrementFloat(key, 3)
		require.NoError(t, err)
		assert.True(t, decremented)

		value, err = redisCache.Cast(key)
		require.NoError(t, err)
		assert.Equal(t, float64(12.3), value.Float64Safe(0))
	})

	t.Run("TTL", func(t *testing.T) {
		key := "redisTTLKey"
		value := "redisTTLValue"
		ttl := 2 * time.Second

		err := redisCache.Put(key, value, &ttl)
		require.NoError(t, err)

		retrievedTTL, err := redisCache.TTL(key)
		require.NoError(t, err)
		assert.LessOrEqual(t, retrievedTTL.Seconds(), ttl.Seconds())
	})
}
