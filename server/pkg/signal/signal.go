package signal

import (
	"errors"
	"sync"
)

var ErrEmpty = errors.New("no topic found")

type topic struct {
	listeners []chan<- struct{}
	mu        *sync.Mutex
}

type Signal interface {
	Subscribe(id string) (<-chan struct{}, func())
	Publish(id string) error
}

type signal struct {
	listeners *sync.Map
}

func NewSignal() Signal {
	return &signal{
		listeners: new(sync.Map),
	}
}

func (c *signal) Subscribe(id string) (<-chan struct{}, func()) {
	topicInf, _ := c.listeners.LoadOrStore(id, &topic{mu: new(sync.Mutex)})
	t := topicInf.(*topic)
	t.mu.Lock()
	defer t.mu.Unlock()
	ch := make(chan struct{}, 1)
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

func (c *signal) Publish(id string) error {
	topicInf, ok := c.listeners.Load(id)
	if !ok {
		return ErrEmpty
	}
	topic := topicInf.(*topic)
	l := len(topic.listeners)
	if l == 0 {
		return ErrEmpty
	}
	for i := 0; i < l; i++ {
		topic.listeners[i] <- struct{}{}
	}
	return nil
}

func (c *signal) HasListener(id string) bool {
	_, ok := c.listeners.Load(id)
	return ok
}
