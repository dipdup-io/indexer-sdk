# Messages

Package implements PubSub pattern via channels

## Install

```bash
go get  github.com/dipdup-net/indexer-sdk/messages
```

## Package structure

Package has three structures: `Publicher`, `Subscriber` and `Message`. 
Example of code can be found [here](/examples/messages/main.go)

### Publisher

Publisher is the structure that notify all subscribers about events.

```go
publisher := messages.NewPublisher()
```

It has following methods:

```go
// Notify - notifies all subscribers with msg
func (publisher *Publisher) Notify(msg *Message) 

// Subscribe - subscribes `subscriber` to `id`
func (publisher *Publisher) Subscribe(subscriber *Subscriber, id SubscriptionID)

// Unsubscribe - unsubscribes `subscriber` from `id`
func (publisher *Publisher) Unsubscribe(subscriber *Subscriber, id SubscriptionID)
```

### Subscriber

Subscriber is the structure which waits messages from publisher.

```go
subscriber := messages.NewSubscriber()
```

It has following methods:

```go
// Listen - waits new message from publisher
func (s *Subscriber) Listen() <-chan *Message

// Close - function that clears subscriber's state
func (s *Subscriber) Close() error
```

### Message

Message is the data structure. It has subscription id field and data.

```go
msg :=  NewMessage(SubscriptionID(100), "message")
```

It has following methods:

```go
// SubscriptionID - returns message's subscription identity
func (msg *Message) SubscriptionID() SubscriptionID 

// Data - returns message's data
func (msg *Message) Data() any
```