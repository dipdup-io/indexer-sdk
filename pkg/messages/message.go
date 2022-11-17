package messages

// Message -
type Message struct {
	id   SubscriptionID
	data any
}

// NewMessage -
func NewMessage(id SubscriptionID, data any) *Message {
	return &Message{
		id:   id,
		data: data,
	}
}

// SubscriptionID - returns message's subscription identity
func (msg *Message) SubscriptionID() SubscriptionID {
	return msg.id
}

// Data - returns message's data
func (msg *Message) Data() any {
	return msg.data
}
