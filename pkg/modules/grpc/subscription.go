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
	m  map[string]Subscription[T, P]
	mx sync.RWMutex
}

// NewSubscriptions -
func NewSubscriptions[T any, P any]() *Subscriptions[T, P] {
	return &Subscriptions[T, P]{
		m: make(map[string]Subscription[T, P]),
	}
}

// NotifyAll -
func (s *Subscriptions[T, P]) NotifyAll(typ T, msg P) {
	s.mx.RLock()
	{
		for _, sub := range s.m {
			if sub != nil && sub.Filter(typ) {
				sub.Send(msg)
			}
		}
	}
	s.mx.Unlock()
}

// Add -
func (s *Subscriptions[T, P]) Add(id string, subscription Subscription[T, P]) {
	s.mx.Lock()
	{
		s.m[id] = subscription
	}
	s.mx.Unlock()
}

// Remove -
func (s *Subscriptions[T, P]) Remove(id string) error {
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
func (s *Subscriptions[T, P]) Get(id string) (Subscription[T, P], bool) {
	defer s.mx.Unlock()
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
