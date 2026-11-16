// Copyright 2024 eve.  All rights reserved.

package db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisOptions defines optsions for mysql database.
type RedisOptions struct {
	Addr         string
	Username     string
	Password     string
	Database     int
	MaxRetries   int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolTimeout  time.Duration
	PoolSize     int
}

// NewRedis create a new redis client with the given options.
func NewRedis(opts *RedisOptions) (*redis.Client, error) {
	options := &redis.Options{
		Addr:         opts.Addr,
		Username:     opts.Username,
		Password:     opts.Password,
		DB:           opts.Database,
		MaxRetries:   opts.MaxRetries,
		MinIdleConns: opts.MinIdleConns,
		DialTimeout:  opts.DialTimeout,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		PoolTimeout:  opts.PoolTimeout,
		PoolSize:     opts.PoolSize,
	}

	rdb := redis.NewClient(options)

	// check redis if is ok，设置超时避免连接异常时无限阻塞
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		_ = rdb.Close()
		return nil, err
	}

	return rdb, nil
}
