package cache

import "time"

// RateLimiter defines the interface for a rate limiter.
type RateLimiter interface {
	// Hit decrements the user's remaining attempts.
	// Returns an error if the operation fails.
	Hit() error

	// Lock locks the rate limiter, preventing further attempts.
	// Returns an error if the operation fails.
	Lock() error

	// Reset resets the rate limiter to its initial state.
	// Returns an error if the operation fails.
	Reset() error

	// Clear removes the rate limiter from the cache.
	// Returns an error if the operation fails.
	Clear() error

	// MustLock checks if the rate limiter should be locked.
	// Returns a boolean indicating the lock state and an error if the operation fails.
	MustLock() (bool, error)

	// TotalAttempts returns the total number of attempts made by the user.
	// Returns the total attempts and an error if the operation fails.
	TotalAttempts() (uint32, error)

	// RetriesLeft returns the number of remaining attempts for the user.
	// Returns the remaining attempts and an error if the operation fails.
	RetriesLeft() (uint32, error)

	// AvailableIn returns the time duration until the rate limiter unlocks.
	// Returns the time until unlock and an error if the operation fails.
	AvailableIn() (time.Duration, error)
}

// limiter is the concrete implementation of the RateLimiter interface.
type limiter struct {
	name        string
	maxAttempts uint32
	ttl         time.Duration
	cache       Cache
}

// NewRateLimiter creates and returns a new rate limiter instance.
func NewRateLimiter(name string, maxAttempts uint32, ttl time.Duration, cache Cache) RateLimiter {
	return &limiter{
		name:        "limiter " + name,
		maxAttempts: maxAttempts,
		ttl:         ttl,
		cache:       cache,
	}
}

func (l *limiter) Hit() error {
	exists, err := l.cache.Decrement(l.name, 1)
	if err != nil {
		return err
	}

	if !exists {
		return l.cache.Put(l.name, l.maxAttempts-1, &l.ttl)
	}

	return nil
}

func (l *limiter) Lock() error {
	exists, err := l.cache.Update(l.name, 0)
	if err != nil {
		return err
	}

	if !exists {
		return l.cache.Put(l.name, 0, &l.ttl)
	}

	return nil
}

func (l *limiter) Reset() error {
	return l.cache.Put(l.name, l.maxAttempts, &l.ttl)
}

func (l *limiter) Clear() error {
	return l.cache.Forget(l.name)
}

func (l *limiter) MustLock() (bool, error) {
	caster, err := l.cache.Cast(l.name)
	if err != nil {
		return true, err
	}

	if caster.IsNil() {
		return false, nil
	}

	num, err := caster.Int()
	if err != nil {
		return true, err
	}

	return num <= 0, nil
}

func (l *limiter) TotalAttempts() (uint32, error) {
	caster, err := l.cache.Cast(l.name)
	if err != nil {
		return 0, err
	}

	if caster.IsNil() {
		return 0, nil
	}

	num, err := caster.Int()
	if err != nil {
		return 0, err
	}

	if num < int(l.maxAttempts) {
		num = int(l.maxAttempts)
	}

	return l.maxAttempts - uint32(num), nil
}

func (l *limiter) RetriesLeft() (uint32, error) {
	caster, err := l.cache.Cast(l.name)
	if err != nil {
		return 0, err
	}

	if caster.IsNil() {
		return 0, nil
	}

	num, err := caster.Int()
	if err != nil {
		return 0, err
	}

	if num < 0 {
		num = 0
	}

	return uint32(num), nil
}

func (l *limiter) AvailableIn() (time.Duration, error) {
	ttl, err := l.cache.TTL(l.name)
	if err != nil {
		return 0, err
	}

	return ttl, nil
}
