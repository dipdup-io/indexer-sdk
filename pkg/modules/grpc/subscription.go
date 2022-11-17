package grpc

import (
	"io"
	"sync"
)

// Subscription - general interface for subscription
type Subscription[T any, P any] interface {
	Filter(typ T) bool
	Send(msg P)
	Listen() <-chan P
	io.Closer
}

// Subscriptions -
type Subscriptions[T any, P any] struct {
	m  map[uint64]Subscription[T, P]
	mx *sync.RWMutex
}

// NewSubscriptions -
func NewSubscriptions[T any, P any]() *Subscriptions[T, P] {
	return &Subscriptions[T, P]{
		m:  make(map[uint64]Subscription[T, P]),
		mx: new(sync.RWMutex),
	}
}

// NotifyAll -
func (s *Subscriptions[T, P]) NotifyAll(typ T, converter func(uint64, T) P) {
	s.mx.RLock()
	{
		for id, sub := range s.m {
			if sub != nil && sub.Filter(typ) {
				sub.Send(converter(id, typ))
			}
		}
	}
	s.mx.RUnlock()
}

// Add -
func (s *Subscriptions[T, P]) Add(id uint64, subscription Subscription[T, P]) {
	s.mx.Lock()
	{
		s.m[id] = subscription
	}
	s.mx.Unlock()
}

// Remove -
func (s *Subscriptions[T, P]) Remove(id uint64) error {
	s.mx.Lock()
	{
		if subs, ok := s.m[id]; ok {
			if err := subs.Close(); err != nil {
				return err
			}
			delete(s.m, id)
		}
	}
	s.mx.Unlock()
	return nil
}

// Get -
func (s *Subscriptions[T, P]) Get(id uint64) (Subscription[T, P], bool) {
	defer s.mx.RUnlock()
	s.mx.RLock()
	subs, ok := s.m[id]
	return subs, ok
}

// Close -
func (s *Subscriptions[T, P]) Close() error {
	s.mx.Lock()
	defer s.mx.Unlock()

	for _, subs := range s.m {
		if subs == nil {
			continue
		}
		if err := subs.Close(); err != nil {
			return err
		}
	}
	return nil
}
