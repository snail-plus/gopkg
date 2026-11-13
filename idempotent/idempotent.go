// Copyright 2024 eve.  All rights reserved.

package idempotent

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/snail-plus/gopkg/log"
)

// redis lua script(read => delete => get delete flag).
const (
	lua string = `
local current = redis.call('GET', KEYS[1])
if current == false then
    return '-1';
end
local del = redis.call('DEL', KEYS[1])
if del == 1 then
     return '1';
else
     return '0';
end
`
)

type Idempotent struct {
	ops Options
}

func New(options ...func(*Options)) *Idempotent {
	ops := getOptionsOrSetDefault(nil)
	for _, f := range options {
		f(ops)
	}
	return &Idempotent{ops: *ops}
}

func (i *Idempotent) Token(ctx context.Context) string {
	if i.ops.redis == nil {
		log.C(ctx).Warnw("please enable redis, otherwise the idempotent is invalid")
		return ""
	}

	token := uuid.NewString()
	key := fmt.Sprintf("%s_%s", i.ops.prefix, token)
	// 写入失败时返回空 token，避免发放一个注定校验失败的 token
	if err := i.ops.redis.Set(ctx, key, true, time.Duration(i.ops.expire)*time.Minute).Err(); err != nil {
		log.C(ctx).Errorw(err, "idempotent token set failed")
		return ""
	}

	return token
}

func (i *Idempotent) Check(ctx context.Context, token string) bool {
	if i.ops.redis == nil {
		log.C(ctx).Warnw("please enable redis, otherwise the idempotent is invalid")
		return true
	}

	res, err := i.ops.redis.Eval(ctx, lua, []string{fmt.Sprintf("%s_%s", i.ops.prefix, token)}).Result()
	if err != nil || res != "1" {
		return false
	}

	return true
}
