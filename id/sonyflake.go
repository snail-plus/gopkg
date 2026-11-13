// Copyright 2024 eve.  All rights reserved.

package id

import (
	"context"
	"fmt"
	"time"

	"github.com/sony/sonyflake"
)

type Sonyflake struct {
	ops   SonyflakeOptions
	sf    *sonyflake.Sonyflake
	Error error
}

// NewSonyflake can get a unique code by id(You need to ensure that id is unique).
func NewSonyflake(options ...func(*SonyflakeOptions)) *Sonyflake {
	ops := getSonyflakeOptionsOrSetDefault(nil)
	for _, f := range options {
		f(ops)
	}
	sf := &Sonyflake{
		ops: *ops,
	}
	st := sonyflake.Settings{
		StartTime: ops.startTime,
	}
	if ops.machineId > 0 {
		st.MachineID = func() (uint16, error) {
			return ops.machineId, nil
		}
	}
	ins := sonyflake.NewSonyflake(st)
	if ins == nil {
		sf.Error = fmt.Errorf("create sonyflake failed")
		return sf
	}
	// 预热一次，提前暴露 startTime 非法等问题
	if _, err := ins.NextID(); err != nil {
		sf.Error = fmt.Errorf("invalid start time: %w", err)
		return sf
	}
	sf.sf = ins
	return sf
}

// maxRetrySleep 限制 Id 重试时的最大退避间隔，避免时钟回拨时退避时间无限增长.
const maxRetrySleep = 64 * time.Millisecond

func (s *Sonyflake) Id(ctx context.Context) (id uint64) {
	if s.Error != nil || s.sf == nil {
		return
	}
	id, err := s.sf.NextID()
	if err == nil {
		return
	}

	// 指数退避重试，响应 ctx 取消，避免调用方被永久阻塞
	sleep := time.Millisecond
	for {
		timer := time.NewTimer(sleep)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return 0
		}

		sleep = min(sleep*2, maxRetrySleep)
		id, err = s.sf.NextID()
		if err == nil {
			return
		}
	}
}
