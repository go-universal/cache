# Cache Library Documentation

![GitHub Tag](https://img.shields.io/github/v/tag/go-universal/cache?sort=semver&label=version)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-universal/cache.svg)](https://pkg.go.dev/github.com/go-universal/cache)
[![License](https://img.shields.io/badge/license-ISC-blue.svg)](https://github.com/go-universal/cache/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-universal/cache)](https://goreportcard.com/report/github.com/go-universal/cache)
![Contributors](https://img.shields.io/github/contributors/go-universal/cache)
![Issues](https://img.shields.io/github/issues/go-universal/cache)

This library provides a flexible and extensible caching system with support for in-memory and Redis-based caching. It includes utilities for managing rate limiters, queues, and verification codes.

## Cache

The `Cache` interface provides a unified API for caching operations:

- `Put(key string, value any, ttl *time.Duration) error`: Store a value with an optional TTL.
- `Update(key string, value any) (bool, error)`: Update an existing key.
- `PutOrUpdate(key string, value any, ttl *time.Duration) error`: Store or update a value.
- `Get(key string) (any, error)`: Retrieve a value.
- `Pull(key string) (any, error)`: Retrieve and remove a value.
- `Cast(key string) (cast.Caster, error)`: Retrieve and cast a value.
- `Exists(key string) (bool, error)`: Check if a key exists.
- `Forget(key string) error`: Remove a key.
- `TTL(key string) (time.Duration, error)`: Get the TTL of a key.
- `Increment(key string, value int64) (bool, error)`: Increment a numeric value.
- `Decrement(key string, value int64) (bool, error)`: Decrement a numeric value.
- `IncrementFloat(key string, value float64) (bool, error)`: Increment a float value.
- `DecrementFloat(key string, value float64) (bool, error)`: Decrement a float value.

## Memory Cache

The `MemoryCache` is an in-memory implementation of the `Cache` interface:

```go
cache := cache.NewMemoryCache()
ttl := 5 * time.Second
err := cache.Put("key", "value", &ttl)
value, err := cache.Get("key")
```

## Redis Cache

The `RedisCache` is a Redis-based implementation of the `Cache` interface:

```go
redisClient := redis.NewClient(&redis.Options{})
cache := cache.NewRedisCache("prefix", redisClient)
ttl := 5 * time.Second
err := cache.Put("key", "value", &ttl)
value, err := cache.Get("key")
```

## Queue

The `Queue` provides methods for managing a queue:

- `Push(value any) error`: Add a value.
- `Pull() (any, error)`: Retrieve and remove the first item.
- `Pop() (any, error)`: Retrieve and remove the last item.
- `Cast() (cast.Caster, error)`: Retrieve and cast the first item.
- `Length() (int64, error)`: Get the number of items.
- `Clear() error`: Remove all items.

## Rate Limiter

The `RateLimiter` manages rate limits:

- `Hit() error`: Decrement remaining attempts.
- `Lock() error`: Lock the rate limiter.
- `Reset() error`: Reset the rate limiter.
- `Clear() error`: Remove the rate limiter.
- `MustLock() (bool, error)`: Check if locking is required.
- `TotalAttempts() (uint32, error)`: Get total attempts.
- `RetriesLeft() (uint32, error)`: Get remaining attempts.
- `AvailableIn() (time.Duration, error)`: Time until unlock.

## Verification Code

The `VerificationCode` manages verification codes:

- `Set(code string) error`: Store a code.
- `Generate(count uint) (string, error)`: Generate a random code.
- `Clear() error`: Clear the code.
- `Get() (string, error)`: Retrieve the code.
- `Validate(code string) (bool, error)`: Validate a code.
- `Exists() (bool, error)`: Check if a code exists.
- `TTL() (time.Duration, error)`: Get the TTL of the code.

## License

This project is licensed under the ISC License. See the [LICENSE](LICENSE) file for details.
