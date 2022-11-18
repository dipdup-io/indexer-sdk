package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/cron"
)

func main() {
	var cfg Config
	if err := config.Parse("config.yml", &cfg); err != nil {
		log.Panic(err)
	}
	log.Print("config was parsed")

	cronModule, err := cron.NewModule(cfg.Cron)
	if err != nil {
		log.Panic(err)
	}
	log.Print("cron module was created")

	customModule := NewCustomModule()
	log.Print("custom module was created")

	ctx, cancel := context.WithCancel(context.Background())

	cronModule.Subscribe(customModule.Subscriber, "every_second")      // set job name as subscription id
	cronModule.Subscribe(customModule.Subscriber, "every_five_second") // set job name as subscription id
	log.Print("custom module was subscribed on cron module")

	customModule.Start(ctx)
	cronModule.Start(ctx)
	log.Print("modules started")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-signals

	cancel()

	if err := customModule.Close(); err != nil {
		log.Panic(err)
	}
	if err := cronModule.Close(); err != nil {
		log.Panic(err)
	}
}
