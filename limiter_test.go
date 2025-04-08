package cache_test

import (
	"testing"
	"time"

	"github.com/go-universal/cache"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateLimtier(t *testing.T) {
	redisCache := cache.NewRedisCache("test", redis.NewClient(&redis.Options{}))

	// Create a new rate limiter
	name := "test-limiter"
	maxAttempts := uint32(5)
	ttl := 10 * time.Second
	rateLimiter := cache.NewRateLimiter(name, maxAttempts, ttl, redisCache)

	// Test Hit method
	err := rateLimiter.Hit()
	require.NoError(t, err)

	retriesLeft, err := rateLimiter.RetriesLeft()
	require.NoError(t, err)
	assert.Equal(t, maxAttempts-1, retriesLeft)

	// Test Lock method
	err = rateLimiter.Lock()
	require.NoError(t, err)

	mustLock, err := rateLimiter.MustLock()
	require.NoError(t, err)
	assert.True(t, mustLock)

	// Test Reset method
	err = rateLimiter.Reset()
	require.NoError(t, err)

	retriesLeft, err = rateLimiter.RetriesLeft()
	require.NoError(t, err)
	assert.Equal(t, maxAttempts, retriesLeft)

	// Test Clear method
	err = rateLimiter.Clear()
	require.NoError(t, err)

	retriesLeft, err = rateLimiter.RetriesLeft()
	require.NoError(t, err)
	assert.Equal(t, uint32(0), retriesLeft)

	// Test AvailableIn method
	availableIn, err := rateLimiter.AvailableIn()
	require.NoError(t, err)
	assert.LessOrEqual(t, availableIn, ttl)
}
