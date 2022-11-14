package messages

import "sync"

type subscribers map[uint64]*Subscriber

// Publisher -
type Publisher struct {
	subscibers subscribers
	topics     map[Topic]subscribers

	mx sync.RWMutex
}

// NewPublisher -
func NewPublisher() *Publisher {
	return &Publisher{
		subscibers: make(subscribers),
		topics:     make(map[Topic]subscribers),
	}
}

// Notify - notifies all subscribers with msg
func (publisher *Publisher) Notify(msg *Message) {
	if msg == nil {
		return
	}

	defer publisher.mx.RUnlock()
	publisher.mx.RLock()

	if subscribers, ok := publisher.topics[msg.topic]; ok {
		for _, subscriber := range subscribers {
			subscriber.notify(msg)
		}
	}
}

// Subscribe - subscribes `subscriber` to `topic`
func (publisher *Publisher) Subscribe(subscriber *Subscriber, topic Topic) {
	if subscriber == nil {
		return
	}

	defer publisher.mx.Unlock()
	publisher.mx.Lock()

	if _, ok := publisher.topics[topic]; !ok {
		publisher.topics[topic] = make(subscribers)
	}
	subscriber.addTopic(topic)
	publisher.topics[topic][subscriber.id] = subscriber
}

// Unsubscribe - unsubscribes `subscriber` from `topic`
func (publisher *Publisher) Unsubscribe(subscriber *Subscriber, topic Topic) {
	if subscriber == nil {
		return
	}

	defer publisher.mx.Unlock()
	publisher.mx.Lock()

	if subscribers, ok := publisher.topics[topic]; ok {
		delete(subscribers, subscriber.id)
	}

	subscriber.removeTopic(topic)
}
