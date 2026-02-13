package service

import (
	"github.com/nicolascb/eventdoctor/internal/api/response"
)

// orderedMap is a generic map that preserves insertion order.
// It eliminates the boilerplate of maintaining a separate map + key slice
// that appears in aggregation functions.
type orderedMap[K comparable, V any] struct {
	keys   []K
	values map[K]*V
}

func newOrderedMap[K comparable, V any]() orderedMap[K, V] {
	return orderedMap[K, V]{
		values: make(map[K]*V),
	}
}

// getOrCreate returns a pointer to the existing value or creates a new one
// using the init function. Preserves insertion order.
func (m *orderedMap[K, V]) getOrCreate(key K, init func() V) *V {
	if v, ok := m.values[key]; ok {
		return v
	}

	val := init()
	m.values[key] = &val
	m.keys = append(m.keys, key)
	return m.values[key]
}

// collect returns values in insertion order.
func (m *orderedMap[K, V]) collect() []V {
	result := make([]V, 0, len(m.keys))
	for _, key := range m.keys {
		result = append(result, *m.values[key])
	}
	return result
}

// findOrAppend searches for an element in a slice by the match predicate.
// If not found, creates a new one using init and appends it to the slice.
// Returns a pointer to the element (existing or new).
func findOrAppend[T any](slice *[]T, match func(*T) bool, init func() T) *T {
	for i := range *slice {
		if match(&(*slice)[i]) {
			return &(*slice)[i]
		}
	}

	*slice = append(*slice, init())
	return &(*slice)[len(*slice)-1]
}

// appendHeader adds a header to the slice if the name is not nil (LEFT JOIN may return NULL).
func appendHeader(headers *[]response.EventHeaderView, name, description *string) {
	if name == nil {
		return
	}

	desc := ""
	if description != nil {
		desc = *description
	}

	*headers = append(*headers, response.EventHeaderView{
		Name:        *name,
		Description: desc,
	})
}
