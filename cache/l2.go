// Copyright 2024 eve.  All rights reserved.

package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"

	"github.com/snail-plus/gopkg/cache/store"
)

// L2Cache represents a two-level cache configuration.
type L2Cache[T any] struct {
	// Options for enabling/disabling caches
	opts *L2Options
	// Local in-memory cache.
	local *ristretto.Cache
	// Remote cache backend
	remote Cache[T]
}

// NewL2 instantiates a new L2 cache.
func NewL2[T any](remote Cache[T], options ...L2Option) *L2Cache[T] {
	opts := NewL2Options()
	for _, opt := range options {
		opt(opts)
	}

	cfg := &ristretto.Config{}
	opts.ApplyTo(cfg)

	// 构造失败通常是参数非法，在初始化阶段直接暴露问题，避免运行期 nil panic
	local, err := ristretto.NewCache(cfg)
	if err != nil {
		panic(fmt.Errorf("create ristretto local cache failed: %w", err))
	}
	return &L2Cache[T]{
		opts:   opts,
		local:  local,
		remote: remote,
	}
}

// Get returns the obj stored in cache if it exists.
func (c *L2Cache[T]) Get(ctx context.Context, key any) (T, error) {
	value, _, err := c.GetWithTTL(ctx, key)
	return value, err
}

// GetWithTTL returns the obj stored in cache and its corresponding TTL, also a bool that is true if the
// item was found and is not expired.
func (c *L2Cache[T]) GetWithTTL(ctx context.Context, key any) (T, time.Duration, error) {
	if !c.opts.Disable {
		ttl, found := c.local.GetTTL(keyFunc(key))
		if !found {
			return *new(T), 0, store.ErrKeyNotFound
		}

		raw, found := c.local.Get(keyFunc(key))
		if !found {
			return *new(T), 0, store.ErrKeyNotFound
		}

		// ristretto 存储的是 any，类型不匹配时返回错误而非 panic
		value, ok := raw.(T)
		if !ok {
			return *new(T), 0, fmt.Errorf("local cache value type mismatch: got %T", raw)
		}

		return value, ttl, nil
	}

	return c.remote.GetWithTTL(ctx, key)
}

// Set populates the cache item using the given key.
func (c *L2Cache[T]) Set(ctx context.Context, key any, obj T) error {
	if !c.opts.Disable {
		_ = c.local.Set(keyFunc(key), obj, 0)
	}

	return c.remote.Set(ctx, key, obj)
}

// SetWithTTL populates the cache item using the given key and TTL.
func (c *L2Cache[T]) SetWithTTL(ctx context.Context, key any, obj T, ttl time.Duration) error {
	if !c.opts.Disable {
		_ = c.local.SetWithTTL(keyFunc(key), obj, 0, ttl)
	}

	return c.remote.SetWithTTL(ctx, key, obj, ttl)
}

// Del removes the cache item using the given key.
func (c *L2Cache[T]) Del(ctx context.Context, key any) error {
	if !c.opts.Disable {
		c.local.Del(keyFunc(key))
	}
	return c.remote.Del(ctx, key)
}

// Clear resets all cache data.
func (c *L2Cache[T]) Clear(ctx context.Context) error {
	if !c.opts.Disable {
		c.local.Clear()
	}
	return c.remote.Clear(ctx)
}

// // Wait waits for all cache operations to complete.
func (c *L2Cache[T]) Wait(ctx context.Context) {
	if !c.opts.Disable {
		c.local.Wait()
	}
	c.remote.Wait(ctx)
}
