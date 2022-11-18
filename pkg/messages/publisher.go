package messages

import (
	"sync"
)

type subscribers map[SubscriptionID]*Subscriber

// Publisher -
type Publisher struct {
	subscribers subscribers

	mx sync.RWMutex
}

// NewPublisher -
func NewPublisher() *Publisher {
	return &Publisher{
		subscribers: make(subscribers),
	}
}

// Notify - notifies all subscribers with msg
func (publisher *Publisher) Notify(msg *Message) {
	if msg == nil {
		return
	}

	defer publisher.mx.RUnlock()
	publisher.mx.RLock()

	if subscriber, ok := publisher.subscribers[msg.id]; ok {
		subscriber.notify(msg)
	}
}

// Subscribe - subscribes `subscriber` to `id`
func (publisher *Publisher) Subscribe(subscriber *Subscriber, id SubscriptionID) {
	if subscriber == nil {
		return
	}

	defer publisher.mx.Unlock()
	publisher.mx.Lock()

	if _, ok := publisher.subscribers[id]; !ok {
		publisher.subscribers[id] = subscriber
	}
	subscriber.addTopic(id)
}

// Unsubscribe - unsubscribes `subscriber` from `id`
func (publisher *Publisher) Unsubscribe(subscriber *Subscriber, id SubscriptionID) {
	if subscriber == nil {
		return
	}

	defer publisher.mx.Unlock()
	publisher.mx.Lock()

	delete(publisher.subscribers, id)
	subscriber.removeTopic(id)
}
