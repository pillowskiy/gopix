package signal

import (
	"errors"
	"sync"
)

var ErrEmpty = errors.New("no topic found")

type topic[T any] struct {
	listeners []chan<- T
	mu        *sync.Mutex
}

type Signal[T any] struct {
	listeners *sync.Map
}

func NewSignal[T any]() *Signal[T] {
	return &Signal[T]{
		listeners: new(sync.Map),
	}
}

func (c *Signal[T]) Subscribe(id string) (<-chan T, func()) {
	topicInf, _ := c.listeners.LoadOrStore(id, &topic[T]{mu: new(sync.Mutex)})
	t := topicInf.(*topic[T])
	t.mu.Lock()
	defer t.mu.Unlock()
	ch := make(chan T, 1)
	t.listeners = append(t.listeners, ch)
	return ch, func() {
		t.mu.Lock()
		defer t.mu.Unlock()
		for i := 0; i < len(t.listeners); i++ {
			if t.listeners[i] == ch {
				t.listeners = append(t.listeners[:i], t.listeners[i+1:]...)
			}
		}
	}
}

func (c *Signal[T]) Publish(id string, item T) error {
	topicInf, ok := c.listeners.Load(id)
	if !ok {
		return ErrEmpty
	}
	topic := topicInf.(*topic[T])
	l := len(topic.listeners)
	if l == 0 {
		return ErrEmpty
	}
	for i := 0; i < l; i++ {
		topic.listeners[i] <- item
	}
	return nil
}
