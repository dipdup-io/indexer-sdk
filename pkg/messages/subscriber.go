package messages

import (
	"sync"
)

// SubscriptionID -
type SubscriptionID any

// Subscriber -
type Subscriber struct {
	subscriptions map[SubscriptionID]struct{}
	messages      chan *Message

	mx sync.RWMutex
}

// NewSubscriber -
func NewSubscriber() *Subscriber {
	return &Subscriber{
		subscriptions: make(map[SubscriptionID]struct{}),
		messages:      make(chan *Message, 1024),
	}
}

func (s *Subscriber) addTopic(id SubscriptionID) {
	defer s.mx.Unlock()
	s.mx.Lock()

	s.subscriptions[id] = struct{}{}
}

func (s *Subscriber) removeTopic(id SubscriptionID) {
	defer s.mx.Unlock()
	s.mx.Lock()

	delete(s.subscriptions, id)
}

// Listen - waits new message from publisher
func (s *Subscriber) Listen() <-chan *Message {
	return s.messages
}

func (s *Subscriber) notify(msg *Message) {
	select {
	case s.messages <- msg:
	default:
	}
}

// Close - function that clears subscriber's state
func (s *Subscriber) Close() error {
	close(s.messages)
	return nil
}
