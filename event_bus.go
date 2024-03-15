package golaze

import (
	"sync"
)

type EventBusConfig struct {
	Subscribers map[string][]EventBusHandler
	Shutdown    chan bool
	lock        *sync.RWMutex
}

type EventBusHandler interface {
	Handle(event *Event) error
}

type Event struct {
	Data interface{}
}

type EventBus struct {
	*EventBusConfig
}

func NewEventBus(config *EventBusConfig) *EventBus {
	if config.Subscribers == nil {
		config.Subscribers = make(map[string][]EventBusHandler)
	}

	if config.lock == nil {
		config.lock = &sync.RWMutex{}
	}

	return &EventBus{
		config,
	}
}

func (eb *EventBus) Subscribe(eventType string, handler EventBusHandler) {
	eb.lock.Lock()
	defer eb.lock.Unlock()

	eb.Subscribers[eventType] = append(eb.Subscribers[eventType], handler)
}

func (eb *EventBus) Unsubscribe(eventType string, handler EventBusHandler) {
	eb.lock.Lock()
	defer eb.lock.Unlock()

	if _, ok := eb.Subscribers[eventType]; !ok {
		return
	}

	// Find the handler and remove it from the list
	for i, h := range eb.Subscribers[eventType] {
		if h == handler {
			eb.Subscribers[eventType] = append(eb.Subscribers[eventType][:i], eb.Subscribers[eventType][i+1:]...)
			return
		}
	}

}

func (eb *EventBus) Publish(eventType string, event *Event) {
	eb.lock.RLock()
	defer eb.lock.RUnlock()

	for _, handler := range eb.Subscribers[eventType] {
		go func(handler EventBusHandler) {
			handler.Handle(event)
		}(handler)
	}
}
