package eventbus

import (
	"sync"

	"github.com/snail-plus/gopkg/log"
)

// EventBus 是一个支持多订阅者的、类型安全的内部事件总线
// 注意：结构体本身是非泛型的
type EventBus struct {
	mu     sync.RWMutex
	topics map[string][]any // key: topic string, value: []chan T
}

// New 创建一个新的EventBus实例
func New() *EventBus {
	return &EventBus{
		topics: make(map[string][]any),
	}
}

// --- 公共的、类型安全的泛型函数 ---

// Subscribe 是一个顶层泛型函数，用于订阅事件
func Subscribe[T any](eb *EventBus, topic string) <-chan T {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	subscriberCh := make(chan T, 100)
	eb.topics[topic] = append(eb.topics[topic], subscriberCh)

	return subscriberCh
}

// Unsubscribe 是一个顶层泛型函数，用于取消订阅
// 匹配规则：同类型且底层 channel 一致。返回是否成功移除。
// 注意：不会关闭 channel，避免与 Publish 竞争导致向已关闭的 channel 发送；
// 取消后不再有新事件写入，订阅方可通过 context 或轮询结束 range。
func Unsubscribe[T any](eb *EventBus, topic string, ch <-chan T) bool {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	subs := eb.topics[topic]
	for i, subscriberAny := range subs {
		subscriber, ok := subscriberAny.(chan T)
		if !ok || subscriber != ch {
			continue
		}
		// 复制而非原地删除，避免并发 Publish 持有的切片快照被修改
		newSubs := make([]any, 0, len(subs)-1)
		newSubs = append(newSubs, subs[:i]...)
		newSubs = append(newSubs, subs[i+1:]...)
		eb.topics[topic] = newSubs
		return true
	}
	return false
}

// Publish 是一个顶层泛型函数，用于发布事件
func Publish[T any](eb *EventBus, topic string, event T) {
	// 持有快照后立即释放读锁；Subscribe 的 append 不会影响已取出的切片
	eb.mu.RLock()
	subscribers := eb.topics[topic]
	eb.mu.RUnlock()

	for _, subscriberAny := range subscribers {
		// 同一 topic 被不同类型订阅时跳过，避免类型断言 panic
		subscriber, ok := subscriberAny.(chan T)
		if !ok {
			continue
		}

		select {
		case subscriber <- event:
		default:
			// 订阅者消费过慢导致缓冲满时丢弃事件，保证发布方不被阻塞
			log.Warnw("event dropped because subscriber channel is full", "topic", topic)
		}
	}
}
