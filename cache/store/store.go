// Copyright 2024 eve.  All rights reserved.

package store

import (
	"context"
	"errors"
	"time"
)

var ErrKeyNotFound = errors.New("key not found")

// Store is the interface for all available stores.
type Store interface {
	Get(ctx context.Context, key any) (any, error)
	GetWithTTL(ctx context.Context, key any) (any, time.Duration, error)
	Set(ctx context.Context, key any, value any) error
	SetWithTTL(ctx context.Context, key any, value any, ttl time.Duration) error
	Del(ctx context.Context, key any) error
	// Clear resets all data in the store.
	// 注意：对于 Redis 等共享存储，这会清空整个实例，请谨慎使用。
	Clear(ctx context.Context) error
	Wait(ctx context.Context)
}
