package cache

import "github.com/go-universal/cast"

// Queue represents a thread-safe, nil-safe queue interface.
// It provides methods to perform common queue operations.
type Queue interface {
	// Push adds a value to the end of the queue.
	// Returns an error if the operation fails.
	Push(value any) error

	// Pull retrieves and removes the first item from the queue.
	// Returns the value and an error if the operation fails.
	Pull() (any, error)

	// Pop retrieves and removes the last item from the queue.
	// Returns the value and an error if the operation fails.
	Pop() (any, error)

	// Cast retrieves and removes the first item from the queue,
	// casting it to a `cast.Caster` type.
	// Returns the `Caster` instance and an error if the operation fails.
	Cast() (cast.Caster, error)

	// Length returns the current number of items in the queue.
	// Returns the queue length and an error if the operation fails.
	Length() (int64, error)

	// Clear removes all items from the queue.
	Clear() error
}
