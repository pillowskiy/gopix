package batch

import (
	"fmt"
	"sync"
	"time"
)

type BatchConfig struct {
	Retries int
	MaxSize int
}

type batcher[T interface{}] struct {
	cfg   *BatchConfig
	items []T
	cb    func([]T) error
	mut   sync.Mutex
}

type Batcher[T interface{}] interface {
	Add(T)
	Ticker(tick time.Duration)
	Tick()
}

// Creates new batcher with default config (Max size = 100 and 3 retries)
func New[T interface{}](cb func([]T) error) Batcher[T] {
	return &batcher[T]{cfg: &BatchConfig{Retries: 3, MaxSize: 100}, cb: cb}
}

// Creates new batcher with input config
func NewWithConfig[T interface{}](config *BatchConfig, cb func([]T) error) Batcher[T] {
	return &batcher[T]{cfg: config, cb: cb}
}

// Add item to batcher queue
func (b *batcher[T]) Add(item T) {
	b.mut.Lock()
	b.items = append(b.items, item)
	b.mut.Unlock()

	if len(b.items) >= b.cfg.MaxSize {
		b.Tick()
	}
}

// Ticker for automatic batch processing
func (b *batcher[T]) Ticker(d time.Duration) {
	ticker := time.NewTicker(d)

	defer func() {
		ticker.Stop()
		b.Tick()
	}()

	for {
		select {
		case <-ticker.C:
			b.Tick()
		}
	}
}

// Tick for batch processing
func (b *batcher[T]) Tick() {
	retries := 0

	process := func() error {
		b.mut.Lock()
		defer b.mut.Unlock()
		if err := b.cb(b.items); err != nil {
			return err
		}
		b.items = []T{}
		return nil
	}

	for {
		if err := process(); err != nil {
			retries++
			if retries < b.cfg.Retries {
				continue
			} else {
				fmt.Println(err)
			}
		}

		break
	}
}
