# Cron

Cron is the module which implements a cron specific parser and notifications about scheduler events.

## Usage

Usage of cron module is described by the [example](/examples/cron/).

To import module in your code writing following line:

```go
import "github.com/dipdup-net/indexer-sdk/pkg/modules/cron"
```

Cron module implements interface `Module`. So you can use it like any other module. For example:

```go
// create cron module
cronModule, err := cron.NewModule(cfg.Cron)
if err != nil {
    log.Panic(err)
}
// start cron module
cronModule.Start(ctx)

// your code is here

// close cron module
if err := cronModule.Close(); err != nil {
    log.Panic(err)
}
```

## Config

Default yaml config of cron module contains only one field `jobs`. It's a map of job names to cron pattern. Job names is used like subscription id in inner-message communication.

```yaml
cron:
  jobs:
    half_an_hour: "0 30 * * * *"
    every_minute: "* 1 * * * *"
    every_five_second: "@every 5s"
    every_second: "* * * * * *"
```

## Output

Module sends to its outputs empty struct which notifies all connected modules about scheduled event. Each job of cron module has own output with names pointed in configuration file. So if your module should execute some work on `every_second` scheduled events from example you should connect it:

```go
// with helper function

if err := modules.Connect(cronModule, customModule, "every_second", "every_second"); err != nil {
    log.Panic(err)
}

// or directly to module

if err := cronModule.AttachTo("every_second", customModule.everySecond); err != nil {
    log.Panic(err)
}
```

Example of handling message from cron's outputs:

```go
for {
    select {
    case <-ctx.Done():
        return
    case <-m.everySecond.Listen():
        log.Info().Msg("arrived from cron module")
    case <-m.everyFiveSecond.Listen():
        log.Info().Msg("arrived from cron module")
    }
}
```

`everySecond` and `everyFiveSecond` are inputs of your modules.