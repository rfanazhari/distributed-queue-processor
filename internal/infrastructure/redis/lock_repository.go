package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rfanazhari/distributed-queue-processor/domain/repository"
)

// LockRepository implements the repository.LockRepository interface using Redis
type LockRepository struct {
	client *redis.Client
}

// NewLockRepository creates a new Redis lock repository
func NewLockRepository(client *redis.Client) repository.LockRepository {
	return &LockRepository{
		client: client,
	}
}

// SetLock attempts to set a lock with the given key and TTL
// It uses Redis SETNX command to ensure atomicity
func (r *LockRepository) SetLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	// Use SETNX to set the key only if it doesn't exist
	result, err := r.client.SetNX(ctx, key, "locked", ttl).Result()
	if err != nil {
		return false, err
	}

	return result, nil
}

// ReleaseLock releases a lock with the given key
func (r *LockRepository) ReleaseLock(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// If result is 1, the key was deleted successfully
	return result == 1, nil
}
