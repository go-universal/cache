package cache

import (
	"errors"
	"math"
	"sync"
	"time"

	"github.com/go-universal/cast"
)

// memRecord represents a single cache entry with its data and optional expiry time.
type memRecord struct {
	data   any
	expiry *time.Time
}

// memCache is an in-memory cache implementation with thread-safe operations.
type memCache struct {
	data  map[string]memRecord
	mutex sync.RWMutex
}

// NewMemoryCache creates and returns a new in-memory cache instance.
func NewMemoryCache() Cache {
	return &memCache{
		data: make(map[string]memRecord),
	}
}

func (m *memCache) Put(key string, value any, ttl *time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var expiry *time.Time
	if ttl != nil {
		exp := time.Now().Add(*ttl)
		expiry = &exp
	}

	m.data[key] = memRecord{
		data:   value,
		expiry: expiry,
	}
	return nil
}

func (m *memCache) Update(key string, value any) (bool, error) {
	record, exists := m.read(key)
	if !exists {
		return false, nil
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	record.data = value
	m.data[key] = *record
	return true, nil
}

func (m *memCache) PutOrUpdate(key string, value any, ttl *time.Duration) error {
	ok, err := m.Update(key, value)
	if err != nil {
		return err
	}

	if !ok {
		return m.Put(key, value, ttl)
	}

	return nil
}

func (m *memCache) Get(key string) (any, error) {
	record, exists := m.read(key)
	if !exists {
		return nil, nil
	}

	return record.data, nil
}

func (m *memCache) Pull(key string) (any, error) {
	val, err := m.Get(key)
	if err != nil {
		return nil, err
	}

	if err := m.Forget(key); err != nil {
		return nil, err
	}

	return val, nil
}

func (m *memCache) Cast(key string) (cast.Caster, error) {
	val, err := m.Get(key)
	return cast.NewCaster(val), err
}

func (m *memCache) Exists(key string) (bool, error) {
	_, exists := m.read(key)
	return exists, nil
}

func (m *memCache) Forget(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.data, key)
	return nil
}

func (m *memCache) TTL(key string) (time.Duration, error) {
	record, exists := m.read(key)
	if !exists {
		return 0, nil
	}

	if record.expiry == nil {
		return time.Duration(math.MaxInt64), nil
	}

	return time.Until(*record.expiry), nil
}

func (m *memCache) Increment(key string, value int64) (bool, error) {
	return m.modifyNumericValue(key, value, func(a, b int64) int64 { return a + b })
}

func (m *memCache) Decrement(key string, value int64) (bool, error) {
	return m.modifyNumericValue(key, value, func(a, b int64) int64 { return a - b })
}

func (m *memCache) IncrementFloat(key string, value float64) (bool, error) {
	return m.modifyFloatValue(key, value, func(a, b float64) float64 { return a + b })
}

func (m *memCache) DecrementFloat(key string, value float64) (bool, error) {
	return m.modifyFloatValue(key, value, func(a, b float64) float64 { return a - b })
}

// read retrieves a cache entry by key, ensuring thread safety and handling expiry.
func (m *memCache) read(key string) (*memRecord, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	val, ok := m.data[key]
	if !ok {
		return nil, false
	}

	// Remove expired entries
	if val.expiry != nil && val.expiry.Before(time.Now()) {
		m.mutex.RUnlock()
		m.mutex.Lock()
		delete(m.data, key)
		m.mutex.Unlock()
		m.mutex.RLock()
		return nil, false
	}

	return &val, true
}

// modifyNumericValue is a helper function to modify integer values in the cache.
func (m *memCache) modifyNumericValue(key string, value int64, op func(int64, int64) int64) (bool, error) {
	record, exists := m.read(key)
	if !exists {
		return false, nil
	}

	caster := cast.NewCaster(record.data)
	num, err := caster.Int64()
	if err != nil {
		return false, errors.New("value is not numeric")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	record.data = op(num, value)
	m.data[key] = *record
	return true, nil
}

// modifyFloatValue is a helper function to modify float values in the cache.
func (m *memCache) modifyFloatValue(key string, value float64, op func(float64, float64) float64) (bool, error) {
	record, exists := m.read(key)
	if !exists {
		return false, nil
	}

	caster := cast.NewCaster(record.data)
	num, err := caster.Float64()
	if err != nil {
		return false, errors.New("value is not numeric")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	record.data = op(num, value)
	m.data[key] = *record
	return true, nil
}
