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
