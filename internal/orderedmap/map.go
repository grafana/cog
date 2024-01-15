package orderedmap

import (
	"sort"
)

type Pair[K, V any] struct {
	Key   K
	Value V
}

type Map[K comparable, V any] struct {
	records map[K]V
	order   []K
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		records: make(map[K]V),
	}
}

func (orderedMap *Map[K, V]) Set(key K, value V) {
	if _, found := orderedMap.records[key]; !found {
		orderedMap.order = append(orderedMap.order, key)
	}

	orderedMap.records[key] = value
}

func (orderedMap *Map[K, V]) Get(key K) V {
	return orderedMap.records[key]
}

func (orderedMap *Map[K, V]) Has(key K) bool {
	_, exists := orderedMap.records[key]
	return exists
}

func (orderedMap *Map[K, V]) Len() int {
	return len(orderedMap.order)
}

func (orderedMap *Map[K, V]) Iterate(callback func(key K, value V)) {
	for _, key := range orderedMap.order {
		callback(key, orderedMap.records[key])
	}
}

func FromMap[K string, V any](original map[K]V) *Map[K, V] {
	orderedMap := New[K, V]()

	keys := make([]K, 0, len(original))
	for key := range original {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	for _, k := range keys {
		orderedMap.Set(k, original[k])
	}

	return orderedMap
}
