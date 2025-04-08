package cache

import (
	"context"
	"errors"
	"time"

	"github.com/go-universal/cast"
	"github.com/redis/go-redis/v9"
)

// redisCache is a Redis-based implementation of the Cache interface.
type redisCache struct {
	prefix string
	client *redis.Client
}

// NewRedisCache creates a new Redis cache instance with a given prefix and Redis client.
func NewRedisCache(prefix string, client *redis.Client) Cache {
	return &redisCache{
		prefix: prefix,
		client: client,
	}
}

func (r *redisCache) Put(key string, value any, ttl *time.Duration) error {
	return r.client.Set(
		context.Background(),
		r.prefixer(key),
		value,
		safeValue(ttl, 0),
	).Err()
}

func (r *redisCache) Update(key string, value any) (bool, error) {
	exists, err := r.Exists(key)
	if err != nil || !exists {
		return false, err
	}

	err = r.client.Set(
		context.Background(),
		r.prefixer(key),
		value,
		redis.KeepTTL,
	).Err()
	return err == nil, err
}

func (r *redisCache) PutOrUpdate(key string, value any, ttl *time.Duration) error {
	ok, err := r.Update(key, value)
	if err != nil {
		return err
	}

	if !ok {
		return r.Put(key, value, ttl)
	}

	return nil
}

func (r *redisCache) Get(key string) (any, error) {
	val, err := r.client.Get(
		context.TODO(),
		r.prefixer(key),
	).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	return val, err
}

func (r *redisCache) Pull(key string) (any, error) {
	val, err := r.Get(key)
	if err != nil {
		return nil, err
	}

	if err := r.Forget(key); err != nil {
		return nil, err
	}

	return val, nil
}

func (r *redisCache) Cast(key string) (cast.Caster, error) {
	val, err := r.Get(key)
	if err != nil {
		return nil, err
	}

	return cast.NewCaster(val), nil
}

func (r *redisCache) Exists(key string) (bool, error) {
	exists, err := r.client.Exists(
		context.TODO(),
		r.prefixer(key),
	).Result()

	if errors.Is(err, redis.Nil) {
		return false, nil
	}

	return exists > 0, err
}

func (r *redisCache) Forget(key string) error {
	err := r.client.Del(
		context.TODO(),
		r.prefixer(key),
	).Err()

	if errors.Is(err, redis.Nil) {
		return nil
	}

	return err
}

func (r *redisCache) TTL(key string) (time.Duration, error) {
	ttl, err := r.client.TTL(
		context.Background(),
		r.prefixer(key),
	).Result()

	if errors.Is(err, redis.Nil) {
		return 0, nil
	}

	return ttl, err
}

func (r *redisCache) Increment(key string, value int64) (bool, error) {
	exists, err := r.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = r.client.IncrBy(
		context.Background(),
		r.prefixer(key),
		value,
	).Err()
	return err == nil, err
}

func (r *redisCache) Decrement(key string, value int64) (bool, error) {
	exists, err := r.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = r.client.DecrBy(
		context.Background(),
		r.prefixer(key),
		value,
	).Err()
	return err == nil, err
}

func (r *redisCache) IncrementFloat(key string, value float64) (bool, error) {
	exists, err := r.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = r.client.IncrByFloat(
		context.Background(),
		r.prefixer(key),
		value,
	).Err()
	return err == nil, err
}

func (r *redisCache) DecrementFloat(key string, value float64) (bool, error) {
	exists, err := r.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = r.client.IncrByFloat(
		context.Background(),
		r.prefixer(key),
		-value,
	).Err()
	return err == nil, err
}

// prefixer adds the prefix to a key to create a namespaced key.
func (r *redisCache) prefixer(key string) string {
	return cacheKey(r.prefix, key)
}
