package repository

import (
	"context"
	"time"
)

// LockRepository defines the interface for lock operations
type LockRepository interface {
	// SetLock attempts to set a lock with the given key and TTL
	// Returns true if lock was successfully set, false otherwise
	SetLock(ctx context.Context, key string, ttl time.Duration) (bool, error)

	// ReleaseLock releases a lock with the given key
	// Returns true if lock was successfully released, false otherwise
	ReleaseLock(ctx context.Context, key string) (bool, error)
}
