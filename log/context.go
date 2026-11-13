// Copyright 2024 eve.  All rights reserved.

package log

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const loggerKeyForGin = "OneXLogger"

// contextKey is how we find Loggers in a context.Context.
type contextKey struct{}

// fieldsFromKeyvals 将偶数长度的 keyvals 转换为 zap Field 列表.
func fieldsFromKeyvals(keyvals ...any) []Field {
	keylen := len(keyvals)
	if keylen == 0 || keylen%2 != 0 {
		return nil
	}

	data := make([]zap.Field, 0, keylen/2)
	for i := 0; i < keylen; i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	return data
}

// WithContext returns a copy of context in which the log value is set.
func WithContext(ctx context.Context, keyvals ...any) context.Context {
	l := FromContext(ctx)
	if fields := fieldsFromKeyvals(keyvals...); len(fields) > 0 {
		l = l.With(fields...)
	}

	// In order to be compatible with *gin.Context, the performance loss is acceptable.
	if c, ok := ctx.(*gin.Context); ok {
		c.Set(loggerKeyForGin, l)
		return c
	}

	return context.WithValue(ctx, contextKey{}, l)
}

// FromContext returns a logger with predefined values from a context.Context.
func FromContext(ctx context.Context, keyvals ...any) Logger {
	var key any = contextKey{}
	if _, ok := ctx.(*gin.Context); ok {
		key = loggerKeyForGin
	}

	var log Logger = std
	if ctx != nil {
		if logger, ok := ctx.Value(key).(Logger); ok {
			log = logger
		}
	}

	if fields := fieldsFromKeyvals(keyvals...); len(fields) > 0 {
		return log.With(fields...)
	}

	return log
}

// C represents for `FromContext` with empty keyvals.
func C(ctx context.Context) Logger {
	return FromContext(ctx).AddCallerSkip(-1)
}
