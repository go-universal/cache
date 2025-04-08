package cache

import (
	"context"
	"errors"

	"github.com/go-universal/cast"
	"github.com/redis/go-redis/v9"
)

// redisQueue represents a Redis-backed queue.
type redisQueue struct {
	name   string
	client *redis.Client
}

// NewRedisQueue creates a new Redis queue instance.
func NewRedisQueue(name string, client *redis.Client) Queue {
	return &redisQueue{
		name:   name,
		client: client,
	}
}

func (r *redisQueue) Push(value any) error {
	return r.client.LPush(context.Background(), r.name, value).Err()
}

func (r *redisQueue) Pull() (any, error) {
	val, err := r.client.LPop(context.Background(), r.name).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return val, nil
}

func (r *redisQueue) Pop() (any, error) {
	val, err := r.client.RPop(context.Background(), r.name).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return val, nil
}

func (r *redisQueue) Cast() (cast.Caster, error) {
	val, err := r.Pull()
	return cast.NewCaster(val), err
}

func (r *redisQueue) Length() (int64, error) {
	val, err := r.client.LLen(context.Background(), r.name).Result()

	if errors.Is(err, redis.Nil) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return val, nil
}

func (r *redisQueue) Clear() error {
	return r.client.Del(context.Background(), r.name).Err()
}
