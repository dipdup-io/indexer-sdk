package messages

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"sync"

	"github.com/rs/zerolog/log"
)

// Topic -
type Topic string

// Subscriber -
type Subscriber struct {
	id       uint64
	topics   map[Topic]struct{}
	messages chan *Message

	mx sync.RWMutex
}

// NewSubscriber -
func NewSubscriber() (*Subscriber, error) {
	i, err := id()
	if err != nil {
		return nil, err
	}
	return &Subscriber{
		id:       i,
		topics:   make(map[Topic]struct{}),
		messages: make(chan *Message, 1024),
	}, nil
}

// ID -
func (s *Subscriber) ID() uint64 {
	return s.id
}

func (s *Subscriber) addTopic(topic Topic) {
	defer s.mx.Unlock()
	s.mx.Lock()

	s.topics[topic] = struct{}{}
}

func (s *Subscriber) removeTopic(topic Topic) {
	defer s.mx.Unlock()
	s.mx.Lock()

	delete(s.topics, topic)
}

// Listen - waits new message from publisher
func (s *Subscriber) Listen() <-chan *Message {
	return s.messages
}

func (s *Subscriber) notify(msg *Message) {
	select {
	case s.messages <- msg:
	default:
		log.Warn().Uint64("id", s.id).Msg("can't send message: channel is full")
	}
}

// Close - function that clears subscriber's state
func (s *Subscriber) Close() error {
	close(s.messages)
	return nil
}

func id() (uint64, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.ReadUvarint(bytes.NewBuffer(b))
}
