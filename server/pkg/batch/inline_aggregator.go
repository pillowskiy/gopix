package batch

import "reflect"

type inlineAggregator[T interface{}] struct {
	data []T
}

func NewInlineAggregator[T interface{}]() Aggregator[T] {
	return &inlineAggregator[T]{}
}

func (a *inlineAggregator[T]) CountByGroup(group string) int {
	size := 0
	for _, v := range a.data {
		if a.inferKey(group, v) == group {
			size++
		}
	}
	return size
}

func (a *inlineAggregator[T]) inferKey(group string, item T) string {
	v := reflect.ValueOf(item).Elem()
	return v.FieldByName(group).String()
}

func (a *inlineAggregator[T]) Count() int {
	return len(a.data)
}

func (a *inlineAggregator[T]) Add(item T) {
	a.data = append(a.data, item)
}

func (a *inlineAggregator[T]) Search(group string, cb func(T) bool) *T {
	for _, v := range a.data {
		if cb(v) {
			return &v
		}
	}

	return nil
}

func (a *inlineAggregator[T]) Clear() {
	a.data = []T{}
}

func (a *inlineAggregator[T]) Aggregate() []T {
	itemsCopy := make([]T, len(a.data))
	copy(itemsCopy, a.data)
	return a.data
}
