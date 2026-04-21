package permission

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type PermissionCache interface {
	GetCodeNames(ctx context.Context, userId uuid.UUID) ([]string, error)
	SetCodeNames(ctx context.Context, userId uuid.UUID, codeNames []string) error
	InvalidateUser(ctx context.Context, userId uuid.UUID) error
	InvalidateAll(ctx context.Context) error
}

type permissionCache struct {
	rdb *redis.Client
}

func NewPermissionCache(rdb *redis.Client) PermissionCache {
	return &permissionCache{rdb: rdb}
}

const (
	permCachePrefix = "perm:"
	permCacheTTL    = 5 * time.Minute
)

func permKey(userId uuid.UUID) string {
	return permCachePrefix + userId.String()
}

func (c *permissionCache) GetCodeNames(ctx context.Context, userId uuid.UUID) ([]string, error) {
	return c.rdb.SMembers(ctx, permKey(userId)).Result()
}

func (c *permissionCache) SetCodeNames(ctx context.Context, userId uuid.UUID, codeNames []string) error {
	if len(codeNames) == 0 {
		return nil
	}

	key := permKey(userId)
	pipe := c.rdb.Pipeline()
	pipe.Del(ctx, key)

	members := make([]interface{}, len(codeNames))
	for i, cn := range codeNames {
		members[i] = cn
	}
	pipe.SAdd(ctx, key, members...)
	pipe.Expire(ctx, key, permCacheTTL)

	_, err := pipe.Exec(ctx)
	return err
}

func (c *permissionCache) InvalidateUser(ctx context.Context, userId uuid.UUID) error {
	return c.rdb.Del(ctx, permKey(userId)).Err()
}

func (c *permissionCache) InvalidateAll(ctx context.Context) error {
	var cursor uint64
	var errs []error
	for {
		keys, nextCursor, err := c.rdb.Scan(ctx, cursor, permCachePrefix+"*", 100).Result()
		if err != nil {
			errs = append(errs, err)
			break
		}
		if len(keys) > 0 {
			if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
				log.Printf("failed to delete permission cache keys: %v", err)
				errs = append(errs, err)
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return errors.Join(errs...)
}

// NoOpPermissionCache is a no-op implementation for testing without Redis.
type noOpPermissionCache struct{}

func NewNoOpPermissionCache() PermissionCache {
	return &noOpPermissionCache{}
}

func (c *noOpPermissionCache) GetCodeNames(ctx context.Context, userId uuid.UUID) ([]string, error) {
	return nil, redis.Nil
}

func (c *noOpPermissionCache) SetCodeNames(ctx context.Context, userId uuid.UUID, codeNames []string) error {
	return nil
}

func (c *noOpPermissionCache) InvalidateUser(ctx context.Context, userId uuid.UUID) error {
	return nil
}

func (c *noOpPermissionCache) InvalidateAll(ctx context.Context) error {
	return nil
}
