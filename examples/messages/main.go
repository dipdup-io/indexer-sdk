package main

import (
	"context"
	"log"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/messages"
)

func main() {
	publisher := messages.NewPublisher()

	subscriber := messages.NewSubscriber()

	id := messages.SubscriptionID(100)

	publisher.Subscribe(subscriber, id)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		tckr := time.NewTicker(time.Second)
		defer tckr.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tckr.C:
				publisher.Notify(messages.NewMessage(id, "new message"))
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
