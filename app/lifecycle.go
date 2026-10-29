// Copyright 2024 eve.  All rights reserved.

package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/snail-plus/gopkg/log"
)

// Startable 表示一个阻塞式运行的后台任务。实现方应在 ctx 取消后尽快退出，
// 以便主流程优雅关停。
type Startable interface {
	Start(ctx context.Context)
}

// Delayed 返回一个延迟 duration 后再启动 s 的 Startable。
// 等待期间若 ctx 已取消则直接返回，不再拉起 s，避免关停后仍被启动。
func Delayed(s Startable, duration time.Duration) Startable {
	return delayedStartable{startable: s, delay: duration}
}

type delayedStartable struct {
	startable Startable
	delay     time.Duration
}

func (d delayedStartable) Start(ctx context.Context) {
	timer := time.NewTimer(d.delay)
	defer timer.Stop()
	select {
	case <-timer.C:
	case <-ctx.Done():
		return
	}
	d.startable.Start(ctx)
}

// Group 并发运行一组后台任务，并将 main 作为主入口阻塞执行，统一受 ctx 控制。
type Group struct {
	tasks []Startable
}

// NewGroup 根据给定后台任务创建任务组。
func NewGroup(tasks ...Startable) *Group {
	return &Group{tasks: append([]Startable(nil), tasks...)}
}

// Add 向任务组追加后台任务。
func (g *Group) Add(tasks ...Startable) {
	g.tasks = append(g.tasks, tasks...)
}

// Run 并发启动所有后台任务，随后阻塞执行 main 直到其返回。
// main 通常是 HTTP/RPC 服务等主入口，其返回的 error 作为 Run 的返回值。
func (g *Group) Run(ctx context.Context, main func(ctx context.Context) error) error {
	for _, t := range g.tasks {
		go t.Start(ctx)
	}
	return main(ctx)
}

// SignalContext 返回一个收到 SIGINT/SIGTERM 后延迟 shutdownDelay 再取消的 ctx。
// 内部信号监听运行在独立 goroutine 中，退出前应调用返回的 cancel 释放资源。
func SignalContext(parent context.Context, shutdownDelay time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		log.Infof("received signal: %s, shutting down", sig)
		if shutdownDelay > 0 {
			time.Sleep(shutdownDelay)
		}
		cancel()
	}()

	return ctx, cancel
}