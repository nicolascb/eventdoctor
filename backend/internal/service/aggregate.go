package service

import (
	"github.com/nicolascb/eventdoctor/internal/api/response"
)

// orderedMap é um map genérico que preserva a ordem de inserção.
// Elimina o boilerplate de manter um map + slice de keys separados
// que aparece nas funções de agregação.
type orderedMap[K comparable, V any] struct {
	keys   []K
	values map[K]*V
}

func newOrderedMap[K comparable, V any]() orderedMap[K, V] {
	return orderedMap[K, V]{
		values: make(map[K]*V),
	}
}

// getOrCreate retorna o ponteiro para o valor existente ou cria um novo
// usando a função init. Preserva a ordem de inserção.
func (m *orderedMap[K, V]) getOrCreate(key K, init func() V) *V {
	if v, ok := m.values[key]; ok {
		return v
	}

	val := init()
	m.values[key] = &val
	m.keys = append(m.keys, key)
	return m.values[key]
}

// collect retorna os valores na ordem de inserção.
func (m *orderedMap[K, V]) collect() []V {
	result := make([]V, 0, len(m.keys))
	for _, key := range m.keys {
		result = append(result, *m.values[key])
	}
	return result
}

// findOrAppend procura um elemento em um slice pelo predicado match.
// Se não encontrar, cria um novo usando init e o adiciona ao slice.
// Retorna o ponteiro para o elemento (existente ou novo).
func findOrAppend[T any](slice *[]T, match func(*T) bool, init func() T) *T {
	for i := range *slice {
		if match(&(*slice)[i]) {
			return &(*slice)[i]
		}
	}

	*slice = append(*slice, init())
	return &(*slice)[len(*slice)-1]
}

// appendHeader adiciona um header ao slice se o nome não for nil (LEFT JOIN pode trazer NULL).
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
