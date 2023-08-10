package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/cron"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/duplicator"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/printer"
)

func main() {
	cronModule1, err := cron.NewModule(&cron.Config{
		Jobs: map[string]string{
			"ticker": "@every 5s",
		},
	})
	if err != nil {
		log.Panic(err)
	}
	cronModule2, err := cron.NewModule(&cron.Config{
		Jobs: map[string]string{
			"ticker": "* * * * * *",
		},
	})
	if err != nil {
		log.Panic(err)
	}

	dup := duplicator.NewDuplicator(2, 1)

	print := printer.NewPrinter()

	if err := modules.Connect(cronModule1, dup, "ticker", duplicator.GetInputName(0)); err != nil {
		log.Panic(err)
	}
	if err := modules.Connect(cronModule2, dup, "ticker", duplicator.GetInputName(1)); err != nil {
		log.Panic(err)
	}
	if err := modules.Connect(dup, print, duplicator.GetOutputName(0), printer.InputName); err != nil {
		log.Panic(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	print.Start(ctx)
	dup.Start(ctx)
	cronModule2.Start(ctx)
	cronModule1.Start(ctx)

	<-ctx.Done()
	cancel()

	if err := cronModule1.Close(); err != nil {
		log.Print(err)
	}
	if err := cronModule2.Close(); err != nil {
		log.Print(err)
	}
	if err := dup.Close(); err != nil {
		log.Print(err)
	}
	if err := print.Close(); err != nil {
		log.Print(err)
	}
}
