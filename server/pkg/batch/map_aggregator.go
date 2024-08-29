package batch

type mapAggregator[T GroupItem] struct {
	data      map[string][]T
	collected int
}

func NewMapAggregator[T GroupItem]() Aggregator[T] {
	return &mapAggregator[T]{
		data: make(map[string][]T),
	}
}

func (a *mapAggregator[T]) CountByGroup(group string) int {
	return len(a.data[group])
}

func (a *mapAggregator[T]) Count() int {
	return a.collected
}

func (a *mapAggregator[T]) Add(item T) {
	a.data[item.Group()] = append(a.data[item.Group()], item)
	a.collected += 1
}

func (a *mapAggregator[T]) Search(group string, cb func(T) bool) *T {
	if cb == nil {
		return nil
	}

	items, ok := a.data[group]
	if ok {
		for _, v := range items {
			if cb(v) {
				return &v
			}
		}
	}

	return nil
}

func (a *mapAggregator[T]) Clear() {
	a.data = make(map[string][]T)
	a.collected = 0
}

func (a *mapAggregator[T]) Aggregate() []T {
	result := make([]T, 0, a.collected)
	for _, v := range a.data {
		result = append(result, v...)
	}
	return result
}
