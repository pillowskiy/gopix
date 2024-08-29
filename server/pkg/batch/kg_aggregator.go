package batch

// KGItem is an interface for objects that have a group and a key.
type KGItem interface {
	// Group returns the group to which the item belongs.
	Group() string

	// Key returns a unique key for the item.
	Key() string
}

type kgAggregator[T KGItem] struct {
	items   map[string]T
	counter map[string]int
}

// NewKGAggregator creates a new aggregator for items of type T.
// Use this aggregator to store unique items and count their occurrences in groups.
// Items with the same key will be overwritten, but group counts will not be increased.
func NewKGAggregator[T KGItem]() Aggregator[T] {
	return &kgAggregator[T]{
		items:   make(map[string]T),
		counter: make(map[string]int),
	}
}

func (a *kgAggregator[T]) Count() int {
	return len(a.items)
}

func (a *kgAggregator[T]) Add(item T) {
	if _, exists := a.items[item.Key()]; !exists {
		a.counter[item.Group()] += 1
	}

	a.items[item.Key()] = item
}

func (a *kgAggregator[T]) Clear() {
	a.items = make(map[string]T)
	a.counter = make(map[string]int)
}

// Looks up for item by its key
// The cb func in this aggregator does nothing, so you can use nil
func (a *kgAggregator[T]) Search(key string, cb func(T) bool) *T {
	item, exists := a.items[key]
	if !exists {
		return nil
	}

	return &item
}

func (a *kgAggregator[T]) CountByGroup(group string) int {
	return a.counter[group]
}

func (a *kgAggregator[T]) Aggregate() []T {
	result := make([]T, 0, len(a.items))
	for _, v := range a.items {
		result = append(result, v)
	}
	return result
}
