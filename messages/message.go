package messages

// Message -
type Message struct {
	topic Topic
	data  any
}

// NewMessage -
func NewMessage(topic Topic, data any) *Message {
	return &Message{
		topic: topic,
		data:  data,
	}
}

// Topic - returns message's topic
func (msg *Message) Topic() Topic {
	return msg.topic
}

// Data - returns message's data
func (msg *Message) Data() any {
	return msg.data
}
