package eventbus

import (
	"reflect"
	"sync"
)

// EventBus 是一个支持多订阅者的、类型安全的内部事件总线
// 注意：结构体本身是非泛型的
type EventBus struct {
	mu     sync.RWMutex
	topics map[string]any // key: topic string, value: []chan T
}

// New 创建一个新的EventBus实例
func New() *EventBus {
	return &EventBus{
		topics: make(map[string]any),
	}
}

// --- 公共的、类型安全的泛型函数 ---

// Subscribe 是一个顶层泛型函数，用于订阅事件
func Subscribe[T any](eb *EventBus, topic string) <-chan T {
	topicStr := topic

	eb.mu.Lock()
	defer eb.mu.Unlock()

	subscriberCh := make(chan T, 100)

	if existingSubscribers, exists := eb.topics[topicStr]; !exists {
		// 如果这是第一个订阅者，创建一个新的订阅者列表
		newSubscribers := []chan T{subscriberCh}
		eb.topics[topicStr] = newSubscribers
	} else {
		// 如果已有订阅者，获取现有列表，将新订阅者追加进去
		subscribersSlice := reflect.ValueOf(existingSubscribers)
		newSubscribers := reflect.Append(subscribersSlice, reflect.ValueOf(subscriberCh))
		eb.topics[topicStr] = newSubscribers.Interface()
	}

	return subscriberCh
}

// Publish 是一个顶层泛型函数，用于发布事件
func Publish[T any](eb *EventBus, topic string, event T) {
	topicStr := topic

	eb.mu.RLock()
	defer eb.mu.RUnlock()

	subscribersAny, ok := eb.topics[topicStr]
	if !ok {
		return
	}

	subscribers := reflect.ValueOf(subscribersAny)
	if subscribers.Kind() != reflect.Slice || subscribers.Len() == 0 {
		return
	}

	for i := 0; i < subscribers.Len(); i++ {
		subscriber := subscribers.Index(i).Interface().(chan T)

		select {
		case subscriber <- event:
		default:
			// 丢弃事件以保证系统健壮性
		}
	}
}
