package sync

import "sync"
import "iter"

type Map[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

func NewMap[K comparable, V any]() Map[K, V] {
	return Map[K, V]{
		m: make(map[K]V),
	}
}

func (m *Map[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.m = make(map[K]V)
}

func (m *Map[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.m, key)
}

func (m *Map[K, V]) Load(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.m[key]
	return v, ok
}

func (m *Map[K, V]) Store(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.m[key] = value
}

func (m *Map[K, V]) LoadAndDelete(key K) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	v, ok := m.m[key]
	if ok {
		delete(m.m, key)
	}

	return v, ok
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if actual, ok := m.m[key]; ok {
		return actual, true
	}

	m.m[key] = value
	return value, false
}

func (m *Map[K, V]) Swap(key K, value V) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	old, loaded := m.m[key]
	m.m[key] = value

	return old, loaded
}

func (m *Map[K, V]) Range() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.mu.RLock()
		defer m.mu.RUnlock()

		for k, v := range m.m {
			if !yield(k, v) {
				return
			}
		}
	}
}
