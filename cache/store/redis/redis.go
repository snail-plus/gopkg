// Copyright 2024 eve.  All rights reserved.

package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/snail-plus/gopkg/cache/store"
	"github.com/snail-plus/gopkg/log"
)

const (
	// RedisType represents the storage type as a string value.
	RedisType = "redis"
)

// RedisStore is a store for Redis.
type RedisStore struct {
	client *redis.Client
}

// NewRedis creates a new store to Redis instance(s).
func NewRedis(client *redis.Client) *RedisStore {
	return &RedisStore{
		client: client,
	}
}

// Get returns data stored from a given key.
func (s *RedisStore) Get(ctx context.Context, key any) (any, error) {
	obj, err := s.client.Get(ctx, key.(string)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, store.ErrKeyNotFound
	}
	return obj, err
}

// GetWithTTL returns data stored from a given key and its corresponding TTL.
func (s *RedisStore) GetWithTTL(ctx context.Context, key any) (any, time.Duration, error) {
	k, ok := key.(string)
	if !ok {
		return nil, 0, store.ErrKeyNotFound
	}

	// 使用 pipeline 合并 GET 与 TTL，减少一次网络往返
	pipe := s.client.Pipeline()
	getCmd := pipe.Get(ctx, k)
	ttlCmd := pipe.TTL(ctx, k)
	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return nil, 0, err
	}

	if err := getCmd.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, 0, store.ErrKeyNotFound
		}
		return nil, 0, err
	}

	ttl, err := ttlCmd.Result()
	if err != nil {
		return nil, 0, err
	}

	return getCmd.Val(), ttl, nil
}

// Set defines data in Redis for given key identifier.
func (s *RedisStore) Set(ctx context.Context, key any, value any) error {
	return s.client.Set(ctx, key.(string), value, 0).Err()
}

// SetWithTTL defines data in Redis for given key identifier with a TTL.
func (s *RedisStore) SetWithTTL(ctx context.Context, key any, value any, ttl time.Duration) error {
	return s.client.Set(ctx, key.(string), value, ttl).Err()
}

// Del removes data from Redis for given key identifier.
func (s *RedisStore) Del(ctx context.Context, key any) error {
	_, err := s.client.Del(ctx, key.(string)).Result()
	return err
}

// Clear resets all data in the store.
// 注意：对 Redis 后端这意味着 FLUSHALL，会清空整个 Redis 实例（包括其他业务的数据），
// 请勿在共享的 Redis 实例上调用。
func (s *RedisStore) Clear(ctx context.Context) error {
	log.Warnw("redis store Clear() will execute FLUSHALL, all data in the redis instance will be lost")
	return s.client.FlushAll(ctx).Err()
}

func (s *RedisStore) Wait(ctx context.Context) {
}
