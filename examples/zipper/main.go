package main

import (
	"context"
	"log"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/zipper"
)

func main() {
	zip := zipper.NewModule[int]()

	first := NewCustomModule(10, -1, "first")
	second := NewCustomModule(0, 1, "second")

	if err := modules.Connect(first, zip, "output", zipper.FirstInputName); err != nil {
		log.Panic(err)
	}
	if err := modules.Connect(second, zip, "output", zipper.SecondInputName); err != nil {
		log.Panic(err)
	}

	fakeInput := modules.NewInput("fake")
	if err := zip.AttachTo(zipper.OutputName, fakeInput); err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	zip.Start(ctx)
	second.Start(ctx)
	first.Start(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-fakeInput.Listen():
				if !ok {
					return
				}
				log.Println(msg)
			}
		}
	}()

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
	if err := fakeInput.Close(); err != nil {
		log.Panic(err)
	}
}
