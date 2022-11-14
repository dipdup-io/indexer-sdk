package main

import (
	"context"
	"log"
	"time"

	"github.com/dipdup-net/indexer-sdk/messages"
)

func main() {
	publisher := messages.NewPublisher()

	subscriber, err := messages.NewSubscriber()
	if err != nil {
		log.Panic(err)
	}

	topic := messages.Topic("some_topic")

	publisher.Subscribe(subscriber, topic)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		tckr := time.NewTicker(time.Second)
		defer tckr.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tckr.C:
				publisher.Notify(messages.NewMessage(topic, "new message"))
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-subscriber.Listen():
				text := msg.Data().(string)
				log.Print(text)
			}
		}
	}()

	time.Sleep(time.Second * 10)

	cancel()
}
