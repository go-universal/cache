package cache

import (
	"time"

	"github.com/go-universal/cast"
)

// Cache provides a nil-safe interface for caching operations.
type Cache interface {
	// Put stores a value in the cache with the specified key and optional TTL (time-to-live).
	// If ttl is nil, the value is stored indefinitely.
	// Returns an error if the operation fails.
	Put(key string, value any, ttl *time.Duration) error

	// Update updates the value of an existing key in the cache.
	// Returns true if the key exists, and an error if the operation fails.
	Update(key string, value any) (bool, error)

	// PutOrUpdate stores a value in the cache with the specified key and TTL.
	// If the key exists, the existing TTL is preserved.
	// If ttl is nil, the value is stored indefinitely.
	// Returns an error if the operation fails.
	PutOrUpdate(key string, value any, ttl *time.Duration) error

	// Get retrieves the value associated with the specified key from the cache.
	// Returns the value and an error if the operation fails.
	Get(key string) (any, error)

	// Pull retrieves the value associated with the specified key and removes it from the cache.
	// Returns the value and an error if the operation fails.
	Pull(key string) (any, error)

	// Cast retrieves the value associated with the specified key and casts it to a cast.Caster.
	// Returns the casted value and an error if the operation fails.
	Cast(key string) (cast.Caster, error)

	// Exists checks whether a key exists in the cache.
	// Returns true if the key exists, and an error if the operation fails.
	Exists(key string) (bool, error)

	// Forget removes the value associated with the specified key from the cache.
	// Returns an error if the operation fails.
	Forget(key string) error

	// TTL retrieves the time-to-live (TTL) of the value associated with the specified key.
	// Returns the TTL and an error if the operation fails.
	TTL(key string) (time.Duration, error)

	// Increment increases the integer value of the specified key by the given amount.
	// Returns true if the key exists, and an error if the operation fails.
	Increment(key string, value int64) (bool, error)

	// Decrement decreases the integer value of the specified key by the given amount.
	// Returns true if the key exists, and an error if the operation fails.
	Decrement(key string, value int64) (bool, error)

	// IncrementFloat increases the float value of the specified key by the given amount.
	// Returns true if the key exists, and an error if the operation fails.
	IncrementFloat(key string, value float64) (bool, error)

	// DecrementFloat decreases the float value of the specified key by the given amount.
	// Returns true if the key exists, and an error if the operation fails.
	DecrementFloat(key string, value float64) (bool, error)
}
