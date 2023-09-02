package main

import (
	"context"
	"log"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/printer"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/zipper"
)

func main() {
	zip := zipper.NewModule[int]()
	first := NewCustomModule(10, -1, "first")
	second := NewCustomModule(0, 1, "second")
	printerModule := NewPrinter()

	if err := modules.Register(zip, first, second, printerModule); err != nil {
		log.Panic(err)
	}

	if err := modules.Connect("first", zip.Name(), "Output", zipper.FirstInputName); err != nil {
		log.Panic(err)
	}
	if err := modules.Connect("second", zip.Name(), "Output", zipper.SecondInputName); err != nil {
		log.Panic(err)
	}
	if err := modules.Connect(zip.Name(), printerModule.Name(), zipper.OutputName, printer.InputName); err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	printerModule.Start(ctx)
	zip.Start(ctx)
	second.Start(ctx)
	first.Start(ctx)

	<-ctx.Done()

	if err := first.Close(); err != nil {
		log.Panic(err)
	}
	if err := second.Close(); err != nil {
		log.Panic(err)
	}
	if err := zip.Close(); err != nil {
		log.Panic(err)
	}
	if err := printerModule.Close(); err != nil {
		log.Panic(err)
	}
}
