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
	cfg *BatchConfig
	agg Aggregator[T]
	cb  func([]T) error
	mut sync.RWMutex
}

type GroupItem interface {
	Group() string
}

type Synchronizer[T interface{}] interface {
	Search(group string, cb func(T) bool) *T
	CountByGroup(group string) int
}

type Aggregator[T interface{}] interface {
	Count() int
	Add(item T)
	Clear()
	Aggregate() []T

	Synchronizer[T]
}

type Batcher[T interface{}] interface {
	Add(T)
	Ticker(tick time.Duration)
	Tick()

	Synchronizer[T]
}

// Creates new batcher with input config
func NewWithConfig[T interface{}](agg Aggregator[T], cb func([]T) error, config *BatchConfig) Batcher[T] {
	return &batcher[T]{agg: agg, cfg: config, cb: cb}
}

// Add item to batcher queue
func (b *batcher[T]) Add(item T) {
	b.mut.Lock()
	b.agg.Add(item)
	b.mut.Unlock()

	if b.agg.Count() >= b.cfg.MaxSize {
		b.Tick()
	}
}

func (b *batcher[T]) Search(group string, cb func(T) bool) *T {
	return b.agg.Search(group, cb)
}

func (b *batcher[T]) CountByGroup(group string) int {
	b.mut.RLock()
	defer b.mut.RUnlock()

	return b.agg.CountByGroup(group)
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

	b.mut.Lock()
	items := b.agg.Aggregate()
	b.agg.Clear()
	b.mut.Unlock()

	process := func() error {
		if err := b.cb(items); err != nil {
			return err
		}
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
