package cache_test

import (
	"testing"
	"time"

	"github.com/go-universal/cache"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerificationCode(t *testing.T) {
	redisCache := cache.NewRedisCache("test", redis.NewClient(&redis.Options{}))

	// Create a new verification code manager
	codeKey := "test"
	codeValue := "123456"
	ttl := 10 * time.Minute
	verificationCode := cache.NewVerification(codeKey, ttl, redisCache)

	// Test Generate method
	_, err := verificationCode.Generate(6)
	require.NoError(t, err)

	// Test Set method
	err = verificationCode.Set(codeValue)
	require.NoError(t, err)

	// Test Validate method (valid code)
	isValid, err := verificationCode.Validate(codeValue)
	require.NoError(t, err)
	assert.True(t, isValid)

	// Test Validate method (invalid code)
	isValid, err = verificationCode.Validate("654321")
	require.NoError(t, err)
	assert.False(t, isValid)

	// Test Exists method
	exists, err := verificationCode.Exists()
	require.NoError(t, err)
	assert.True(t, exists)

	// Test TTL method
	ttlValue, err := verificationCode.TTL()
	require.NoError(t, err)
	assert.Greater(t, ttlValue, 0*time.Second)

	// Test Clear method
	err = verificationCode.Clear()
	require.NoError(t, err)

	// Verify code is cleared
	isValid, err = verificationCode.Validate(codeValue)
	require.NoError(t, err)
	assert.False(t, isValid)
}
