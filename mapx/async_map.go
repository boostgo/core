package mapx

import (
	"sort"
	"sync"
)

type RWLocker interface {
	RLock()
	RUnlock()
}

// AsyncMap thread safe map implementation.
//
// Contain data map defend by [sync.RWMutex].
//
// Getting keys list is cached (no need iteration).
//
// Can get length/size of [AsyncMap]
type AsyncMap[K comparable, V any] struct {
	keys []K
	data map[K]V
	mx   sync.Locker
	rwMX RWLocker
}

// NewAsyncMap creates AsyncMap
func NewAsyncMap[K comparable, V any](size ...int) *AsyncMap[K, V] {
	const defaultMapSize = 10
	mapSize := defaultMapSize
	if len(size) > 0 {
		mapSize = size[0]
	}

	return &AsyncMap[K, V]{
		keys: make([]K, 0, mapSize),
		data: make(map[K]V, mapSize),
		mx:   &sync.Mutex{},
	}
}

func (am *AsyncMap[K, V]) Locker(locker sync.Locker) *AsyncMap[K, V] {
	am.mx = locker
	return am
}

func (am *AsyncMap[K, V]) RWLocker(locker RWLocker) *AsyncMap[K, V] {
	am.rwMX = locker
	return am
}

// Store provided key & value pair
func (am *AsyncMap[K, V]) Store(key K, value V) *AsyncMap[K, V] {
	am.lock()
	defer am.unlock()

	am.data[key] = value
	am.keys = append(am.keys, key)
	return am
}

// Load get value by provided key
func (am *AsyncMap[K, V]) Load(key K) (V, bool) {
	am.lock()
	defer am.unlock()
	v, ok := am.data[key]
	return v, ok
}

// Keys return all keys
func (am *AsyncMap[K, V]) Keys() []K {
	am.lock()
	defer am.unlock()
	return am.keys
}

// Len returns length of map
func (am *AsyncMap[K, V]) Len() int {
	am.lock()
	defer am.unlock()
	return len(am.data)
}

// Delete element by key
func (am *AsyncMap[K, V]) Delete(key K) *AsyncMap[K, V] {
	am.lock()
	defer am.unlock()

	delete(am.data, key)
	am.keys = removeFromSliceWhere(am.keys, func(keyIterator K) bool {
		return key == keyIterator
	})

	return am
}

// Each iterate over all map elements.
// Stops when provided function return false or when all were keys iterated
func (am *AsyncMap[K, V]) Each(fn func(key K, value V) bool) *AsyncMap[K, V] {
	for _, key := range am.Keys() {
		value, ok := am.Load(key)
		if !ok {
			continue
		}

		if !fn(key, value) {
			break
		}
	}

	return am
}

// Map returns inner map
func (am *AsyncMap[K, V]) Map() map[K]V {
	am.lock()
	defer am.unlock()
	return am.data
}

func (am *AsyncMap[K, V]) lock() {
	if am.rwMX != nil {
		am.rwMX.RLock()
		return
	}

	am.mx.Lock()
}

func (am *AsyncMap[K, V]) unlock() {
	if am.rwMX != nil {
		am.rwMX.RUnlock()
		return
	}

	am.mx.Unlock()
}

func removeFromSlice[T any](source []T, index ...int) []T {
	if len(index) == 0 {
		return source
	}

	if len(index) == 1 {
		i := index[0]

		if i < 0 || i >= len(source) {
			return source
		}

		return append(source[:i], source[i+1:]...)
	}

	sort.Ints(index)

	dst := make([]T, 0, len(source))

	prev := 0
	for _, i := range index {
		if i < 0 || i >= len(source) {
			continue
		}

		dst = append(dst, source[prev:i]...)
		prev = i + 1
	}

	return append(dst, source[prev:]...)
}

func indexOfSlice[T any](source []T, fn func(T) bool) int {
	for index, element := range source {
		if fn(element) {
			return index
		}
	}

	return -1
}

func removeFromSliceWhere[T any](source []T, fn func(T) bool) []T {
	index := indexOfSlice(source, fn)
	if index == -1 {
		return source
	}

	return removeFromSlice(source, index)
}
