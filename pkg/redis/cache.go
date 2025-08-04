package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	DeleteRaw(ctx context.Context, fullKey string) error
	ScanKeys(ctx context.Context, pattern string) ([]string, error)
}

type cacheImpl struct {
	client    *redis.Client
	namespace string
}

func NewCache(client *redis.Client, namespace string) Cache {
	return &cacheImpl{
		client:    client,
		namespace: namespace,
	}
}

func (c *cacheImpl) buildKey(key string) string {
	return c.namespace + ":" + key
}

func (c *cacheImpl) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.client.Set(ctx, c.buildKey(key), value, ttl).Err()
}

func (c *cacheImpl) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, c.buildKey(key)).Result()
}

func (c *cacheImpl) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, c.buildKey(key)).Err()
}

func (c *cacheImpl) DeleteRaw(ctx context.Context, fullKey string) error {
	return c.client.Del(ctx, fullKey).Err()
}

func (c *cacheImpl) ScanKeys(ctx context.Context, pattern string) ([]string, error) {
	var cursor uint64
	var keys []string

	prefixedPattern := c.buildKey(pattern)

	for {
		result, newCursor, err := c.client.Scan(ctx, cursor, prefixedPattern, 10).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, result...)
		cursor = newCursor
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}
