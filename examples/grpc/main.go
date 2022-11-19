package main

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	})

	bind := "127.0.0.1:8889"
	serverCfg := grpc.ServerConfig{
		Bind: bind,
	}

	// create server module
	server, err := NewServer(&serverCfg)
	if err != nil {
		log.Panic().Err(err).Msg("creating server error")
		return
	}

	// creating client module
	client := NewClient(bind)

	ctx, cancel := context.WithCancel(context.Background())

	// starting all modules
	server.Start(ctx)

	if err := client.Connect(ctx); err != nil {
		log.Panic().Err(err).Msg("connecting to server error")
		return
	}
	client.Start(ctx)

	// subscribing to time
	subscriptionID, err := client.SubscribeOnTime(ctx)
	if err != nil {
		log.Panic().Err(err).Msg("subscribing error")
		return
	}
	log.Info().Uint64("subscription_id", subscriptionID).Msg("subscribed")

	// creating custom module which receives notification from client and log it to console.
	module := NewCustomModule()

	if err := modules.Connect(client, module, "time", "input"); err != nil {
		log.Panic().Err(err).Msg("module connection error")
		return
	}

	module.Start(ctx)

	time.Sleep(time.Minute)

	if err := client.UnsubscribeFromTime(ctx, subscriptionID); err != nil {
		log.Panic().Err(err).Msg("unsubscribing error")
		return
	}
	log.Info().Msg("unsubscribed")

	cancel()

	// closing all modules
	if err := client.Close(); err != nil {
		log.Panic().Err(err).Msg("closing client error")
	}
	if err := module.Close(); err != nil {
		log.Panic().Err(err).Msg("closing custo module error")
	}
	if err := server.Close(); err != nil {
		log.Panic().Err(err).Msg("closing server error")
	}
}
