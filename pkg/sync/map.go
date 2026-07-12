// Package sync provides concurrency helpers: a thread-safe generic map
// and a limit/offset pagination iterator.
package sync

import (
	"iter"
	"sync"
)

// Map is a generic map safe for concurrent use by multiple goroutines.
// Access is synchronized with sync.RWMutex: readers proceed in parallel,
// writers take exclusive ownership. Use NewMap to create one; the zero
// value is not usable.
type Map[K comparable, V any] struct {
	m  map[K]V
	mx *sync.RWMutex
}

// NewMap returns an initialized empty Map.
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		m:  make(map[K]V),
		mx: new(sync.RWMutex),
	}
}

// Get returns the value stored for key. The second result reports
// whether the key was present.
func (m *Map[K, V]) Get(key K) (V, bool) {
	m.mx.RLock()
	val, ok := m.m[key]
	m.mx.RUnlock()
	return val, ok
}

// Delete removes the entry for key. It is a no-op if the key is absent.
func (m *Map[K, V]) Delete(key K) {
	m.mx.Lock()
	delete(m.m, key)
	m.mx.Unlock()
}

// Set stores value under key, replacing any existing value.
func (m *Map[K, V]) Set(key K, value V) {
	m.mx.Lock()
	m.m[key] = value
	m.mx.Unlock()
}

// Range calls handler for each entry in unspecified order while holding
// the read lock. Iteration stops early when handler returns a non-nil
// error (which is returned) or true as the second result.
//
// The handler must not call methods that take the write lock (Set,
// Delete, Clear) — that would deadlock.
//
// Range (WARN) does not support nested ranges with Delete in them.
func (m *Map[K, V]) Range(handler func(key K, value V) (error, bool)) error {
	if handler == nil {
		return nil
	}
	m.mx.RLock()
	defer m.mx.RUnlock()

	for k, v := range m.m {
		err, br := handler(k, v)
		if err != nil {
			return err
		}
		if br {
			return nil
		}
	}
	return nil
}

// Clear removes all entries from the map.
func (m *Map[K, V]) Clear() {
	m.mx.Lock()
	clear(m.m)
	m.mx.Unlock()
}

// Len returns the number of entries in the map.
func (m *Map[K, V]) Len() int {
	m.mx.RLock()
	defer m.mx.RUnlock()
	return len(m.m)
}

// Values returns a new slice with all values in unspecified order.
func (m *Map[K, V]) Values() []V {
	arr := make([]V, 0, len(m.m))
	m.mx.RLock()
	for _, v := range m.m {
		arr = append(arr, v)
	}
	m.mx.RUnlock()
	return arr
}

// Keys returns a new slice with all keys in unspecified order.
func (m *Map[K, V]) Keys() []K {
	arr := make([]K, 0, len(m.m))
	m.mx.RLock()
	for k := range m.m {
		arr = append(arr, k)
	}
	m.mx.RUnlock()
	return arr
}

// All returns an iterator over key-value pairs in unspecified order.
// The read lock is held for the whole iteration, so the loop body must
// not call methods that take the write lock (Set, Delete, Clear) —
// that would deadlock.
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.mx.RLock()
		defer m.mx.RUnlock()

		for k, v := range m.m {
			if !yield(k, v) {
				return
			}
		}
	}
}

// AllKeys returns an iterator over keys in unspecified order.
// The read lock is held for the whole iteration; see All for the
// locking caveat.
func (m *Map[K, V]) AllKeys() iter.Seq[K] {
	return func(yield func(K) bool) {
		m.mx.RLock()
		defer m.mx.RUnlock()

		for k := range m.m {
			if !yield(k) {
				return
			}
		}
	}
}

// AllValues returns an iterator over values in unspecified order.
// The read lock is held for the whole iteration; see All for the
// locking caveat.
func (m *Map[K, V]) AllValues() iter.Seq[V] {
	return func(yield func(V) bool) {
		m.mx.RLock()
		defer m.mx.RUnlock()

		for _, v := range m.m {
			if !yield(v) {
				return
			}
		}
	}
}
